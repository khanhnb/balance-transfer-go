package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID        = "mychannel"
	orgName          = "Org1"
	orgAdmin         = "Admin"
	ordererOrgName   = "ordererorg"
	ccID             = "mycc"
	configFileName   = "config.yaml"
	username         = "User1"
	secret           = "thisismysecret"
	identityTypeUser = "user"
)

var secretKey = []byte(secret)
var sdk *fabsdk.FabricSDK
var client *channel.Client

func main() {
	sdk = setupSDK(configFileName)
	// client = getClient(sdk, channelID, username, orgName)
	r := mux.NewRouter()
	r.HandleFunc("/", helloRest).Methods("GET")
	r.HandleFunc("/users", login).Methods("POST")
	r.Handle("/channels", authMiddleware(http.HandlerFunc(createChannel))).Methods("POST")
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r))

}

func createChannel(w http.ResponseWriter, r *http.Request) {

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Add("username", claims["username"].(string))
			r.Header.Add("orgName", claims["orgName"].(string))
			next.ServeHTTP(w, r)

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
		}
	})
}

func login(w http.ResponseWriter, r *http.Request) {
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
		tokenString, err := token.SignedString(secretKey)
		fmt.Println(tokenString, err)
		getRegisteredUser(username, orgName)
		w.Write([]byte(getRegisteredUser(username, orgName)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func helloRest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func setupSDK(configFileName string) *fabsdk.FabricSDK {
	var config = config.FromFile(configFileName)
	sdk, err := fabsdk.New(config)
	if err != nil {
		fmt.Printf("failed to create new SDK: %s\n", err)
	}
	return sdk
}

func getClient(sdk *fabsdk.FabricSDK, channelName string, userName string, orgName string) *channel.Client {
	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Printf("failed to create new channel client: %s\n", err)
	}
	return client
}

func queryCC(client *channel.Client) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: "query", Args: [][]byte{[]byte("a")}},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Printf("failed to query funds: %s\n", err)
	}
	return response.Payload
}

func executeCC(client *channel.Client) []byte {
	response, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: [][]byte{[]byte("a"), []byte("b"), []byte("10")}},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Printf("failed to invoke funds: %s\n", err)
	}
	fmt.Println(response)
	return response.Payload
}
