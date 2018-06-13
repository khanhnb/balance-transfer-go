package utils

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

// QueryCC query chaincode
func QueryCC(client *channel.Client, ccID, fcn string, args [][]byte) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Printf("failed to query funds: %s\n", err)
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
