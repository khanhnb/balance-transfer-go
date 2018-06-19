package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/balance-transfer-go/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var hfc utils.FabricSetup

func main() {

	hfc = utils.FabricSetup{
		AdminUser:         "Admin",
		OrdererOrgName:    "ordererorg",
		ConfigFileName:    "config.yaml",
		Secret:            []byte("thisismysecret"),
		IdentityTypeUser:  "user",
		RegistrarUsername: "admin",
		RegistrarPassword: "adminpw",
		ChannelID:         "mychannel",
	}
	hfc.Init()

	r := mux.NewRouter()
	r.HandleFunc("/users", login).Methods("POST")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(queryCC))).Methods("GET")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(invokeCC))).Methods("POST")
	r.Handle("/channels/{channelName}/blocks/{blockID}", authMiddleware(http.HandlerFunc(getBlockByNumber))).Methods("GET")
	r.Handle("/channels/{channelName}/transactions/{transactionID}", authMiddleware(http.HandlerFunc(getTransactionByID))).Methods("GET")
	r.Handle("/channels/{channelName}", authMiddleware(http.HandlerFunc(getChainInfo))).Methods("GET")
	//r.HandleFunc("/channels/{channelName}", getChainInfo).Methods("GET")
	//r.HandleFunc("/channels/{channelName}/transactions/{transactionID}", getTransactionByID).Methods("GET")
	//r.HandleFunc("/channels/{channelName}/blocks/{blockID}", getBlockByNumber).Methods("GET")
	http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, r))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		if tokenString != "" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return hfc.Secret, nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					r.Header.Add("username", claims["username"].(string))
					r.Header.Add("orgName", claims["orgName"].(string))
					next.ServeHTTP(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== LOGIN ==================")
	// define response
	type response struct {
		Success bool
		Message string
		Token   string
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	username := r.Form.Get("username")
	orgName := r.Form.Get("orgName")
	if username != "" && orgName != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"orgName":  orgName,
			"exp":      time.Now().Unix() + 36000,
		})
		tokenString, err := token.SignedString(hfc.Secret)
		fmt.Println(tokenString, err)
		message, success := utils.GetRegisteredUser(username, orgName, hfc.IdentityTypeUser, hfc.MspClient)
		res := response{
			Success: success,
			Message: message,
		}

		if success {
			res.Token = tokenString
		}
		out, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func queryCC(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== QUERY BY CHAINCODE ==================")

	type args struct {
		array []string
	}

	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	vars := mux.Vars(r)
	fcn := r.URL.Query().Get("fcn")
	tmp := args{}
	json.Unmarshal([]byte(r.URL.Query().Get("args")), &tmp.array)

	client := utils.GetClient(hfc.Sdk, vars["channelName"], username, orgName)
	res := utils.QueryCC(client, vars["chaincodeName"], fcn, utils.GetArgs(tmp.array), r.URL.Query().Get("peer"))
	w.Write(res)
}

func invokeCC(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== INVOKE ON CHAINCODE ==================")
	type invokeBody struct {
		Peers []string
		Fcn   string
		Args  []string
	}
	type response struct {
		TxID string
	}
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	decoder := json.NewDecoder(r.Body)
	body := invokeBody{}
	decoder.Decode(&body)
	client := utils.GetClient(hfc.Sdk, vars["channelName"], username, orgName)
	txid := utils.ExecuteCC(client, vars["chaincodeName"], body.Fcn, utils.GetArgs(body.Args), body.Peers)
	res := response{}
	res.TxID = txid
	out, err := json.Marshal(res)
	if err != nil {
		log.Printf("Marshal res error: %s", err.Error())
	}
	w.Write(out)
}

func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== GET BLOCK BY NUMBER ==================")
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	blockID, _ := strconv.ParseUint(vars["blockID"], 10, 64)

	channelContext := hfc.Sdk.ChannelContext(vars["channelName"], fabsdk.WithUser(username), fabsdk.WithOrg(orgName))
	// ledger client
	client, err := ledger.New(channelContext)

	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}

	w.Write(utils.QueryBlockByNumber(client, blockID, r.URL.Query().Get("peer")))
}

func getTransactionByID(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET TRANSACTION BY TRANSACTION_ID ======================")
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	channelContext := hfc.Sdk.ChannelContext(vars["channelName"], fabsdk.WithUser(username), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Printf("Failed to create new ledger client: %s", err)
	}
	res := utils.QueryTransactionByID(client, vars["transactionID"], r.URL.Query().Get("peer"))
	w.Write(res)
}

func getChainInfo(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET CHANNEL INFORMATION ======================")
	vars := mux.Vars(r)
	username := r.Header.Get("username")
	orgName := r.Header.Get("orgName")
	channelContext := hfc.Sdk.ChannelContext(vars["channelName"], fabsdk.WithUser(username), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Printf("Failed to create new ledger client: %s", err)
	}
	res := utils.QueryChainInfo(client, r.URL.Query().Get("peer"))
	w.Write(res)
}
