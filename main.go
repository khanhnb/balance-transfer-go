package main

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "ordererorg"
	ccID           = "mycc"
)

func main() {
	var config = config.FromFile("config.yaml")
	sdk, err := fabsdk.New(config)
	if err != nil {
		fmt.Sprintf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	clientChannelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"), fabsdk.WithOrg("org1"))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Sprintf("Failed to create new channel client: %s", err)
	}
	value := queryCC(client)

	fmt.Sprintf("Result: %s", string(value[:]))
}

func queryCC(client *channel.Client) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: "query", Args: [][]byte{[]byte("a")}},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Errorf("Failed to query funds: %s", err)
	}
	return response.Payload
}
