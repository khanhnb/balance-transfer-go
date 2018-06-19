package utils

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

// ExecuteCC invoke chaincode
func ExecuteCC(client *channel.Client, ccID, fcn string, args [][]byte, endpoints []string) string {
	response, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoints...))
	if err != nil {
		fmt.Printf("failed to invoke funds: %s\n", err)
	}
	return string(response.TransactionID)
}
