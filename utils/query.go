package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// QueryCC query chaincode
func QueryCC(client *channel.Client, ccID, fcn string, args [][]byte, endpoint string) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("failed to query funds: %s\n", err)
	}
	return response.Payload
}

// QueryBlockByNumber query block by block number
func QueryBlockByNumber(client *ledger.Client, blockID uint64, endpoint string) []byte {
	block, err := client.QueryBlock(blockID, ledger.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %s", err)
	}
	if block.Data == nil {
		log.Fatal("QueryBlockByHash block data is nil")
	}
	return processBlock(block)
}

// QueryTransactionByID get transaction by transaction ID
func QueryTransactionByID(client *ledger.Client, txid string, endpoint string) []byte {
	processedTransaction, err := client.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("QueryTransaction error: %s", err.Error())
	}
	transaction := processTransaction(processedTransaction.GetTransactionEnvelope())
	transaction.ValidationCode = uint8(processedTransaction.GetValidationCode())
	transaction.ValidationCodeName = pb.TxValidationCode_name[int32(transaction.ValidationCode)]
	transactionJSON, _ := json.Marshal(transaction)
	transactionJSONString, _ := Prettyprint(transactionJSON)
	return transactionJSONString
}

// QueryChainInfo Query blockchain info
func QueryChainInfo(client *ledger.Client, endpoint string) []byte {
	blockchainInfo, _ := client.QueryInfo(ledger.WithTargetEndpoints(endpoint))
	type chainInfo struct {
		Height            uint64
		CurrentBlockHash  string
		PreviousBlockHash string
	}
	bci := chainInfo{}
	bci.Height = blockchainInfo.BCI.Height
	bci.CurrentBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.CurrentBlockHash)
	bci.PreviousBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.PreviousBlockHash)
	res, _ := json.Marshal(bci)
	return res
}

// QueryInstalledChaincodes query installed chaincode - must call with Admin context
func QueryInstalledChaincodes(client *resmgmt.Client, endpoint string) []byte {
	chaincodeQueryRes, err := client.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstalledChaincodes: %s", err)
	}
	res, _ := json.Marshal(chaincodeQueryRes)
	return res
}

// QueryInstantiatedChaincodes query instantiated chaincode - must call with Admin context
func QueryInstantiatedChaincodes(client *resmgmt.Client, channelName string, endpoint string) []byte {
	chaincodeQueryRes, err := client.QueryInstantiatedChaincodes(channelName, resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstantiatedChaincodes: %s", err)
	}
	res, _ := json.Marshal(chaincodeQueryRes)
	return res
}

// QueryChannels query channels - must call with Admin context
func QueryChannels(client *resmgmt.Client, endpoint string) []byte {
	channelQueryRes, err := client.QueryChannels(resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryChannels: %s", err)
	}
	res, _ := json.Marshal(channelQueryRes)
	return res
}
