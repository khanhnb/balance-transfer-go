package utils

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/test/integration"
)

// GetClient get client from channelName, userName and orgName
func GetClient(sdk *fabsdk.FabricSDK, channelName string, userName string, orgName string) *channel.Client {
	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Printf("failed to create new channel client: %s\n", err)
	}
	return client
}

// GetRegisteredUser get registered user. If user is not enrolled, enroll new user
func GetRegisteredUser(username, orgName, identityTypeUser string, mspClient *msp.Client) (string, bool) {
	testAttributes := []msp.Attribute{
		{
			Name:  integration.GenerateRandomID(),
			Value: fmt.Sprintf("%s:ecert", integration.GenerateRandomID()),
			ECert: true,
		},
		{
			Name:  integration.GenerateRandomID(),
			Value: fmt.Sprintf("%s:ecert", integration.GenerateRandomID()),
			ECert: true,
		},
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:       username,
		Type:       identityTypeUser,
		Attributes: testAttributes,
		// Affiliation is mandatory. "org1" and "org2" are hardcoded as CA defaults
		// See https://github.com/hyperledger/fabric-ca/blob/release/cmd/fabric-ca-server/config.go
		Affiliation: "org1",
	})
	// err if user is already enrolled
	if err == nil {
		// Enroll the new user
		err = mspClient.Enroll(username, msp.WithSecret(enrollmentSecret))
		log.Printf("secret: %s", enrollmentSecret)
		if err != nil {
			log.Printf("enroll %s failed: %v", username, err)
			return "failed " + err.Error(), false
		}
	}
	return username + " enrolled Successfully", true
}

// GetArgs get [][]byte args from string array
func GetArgs(args []string) [][]byte {
	var result [][]byte
	for _, element := range args {
		result = append(result, []byte(element))
	}
	return result
}
