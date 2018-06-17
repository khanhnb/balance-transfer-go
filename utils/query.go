package utils

import (
	"encoding/json"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// QueryCC query chaincode
func QueryCC(client *channel.Client, ccID, fcn string, args [][]byte) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		log.Fatalf("failed to query funds: %s\n", err)
	}
	return response.Payload
}

// QueryBlockByNumber query block by block number
func QueryBlockByNumber(client *ledger.Client, blockID uint64) []byte {
	block, err := client.QueryBlock(blockID)
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %s", err)
	}
	if block.Data == nil {
		log.Fatal("QueryBlockByHash block data is nil")
	}
	return processBlock(block)
}

// QueryTransactionByID get transaction by transaction ID
func QueryTransactionByID(client *ledger.Client, txid string) []byte {
	processedTransaction, err := client.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		log.Fatalf("QueryTransaction error: %s", err.Error())
	}
	transaction := processTransaction(processedTransaction.GetTransactionEnvelope())
	transaction.ValidationCode = uint8(processedTransaction.GetValidationCode())
	transaction.ValidationCodeName = pb.TxValidationCode_name[int32(transaction.ValidationCode)]
	transactionJSON, _ := json.Marshal(transaction)
	transactionJSONString, _ := prettyprint(transactionJSON)
	return transactionJSONString
}
