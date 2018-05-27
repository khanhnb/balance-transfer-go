package utils

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
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
