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
func GetRegisteredUser(username, orgName, secret, identityTypeUser string, sdk *fabsdk.FabricSDK) (string, bool) {
	ctxProvider := sdk.Context(fabsdk.WithOrg(orgName))
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		log.Fatalf("Failed to create msp client: %s", err.Error())
	}
	signingIdentity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		log.Printf("Check if user %s is enrolled: %s", username, err.Error())
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
		identity, err := mspClient.GetIdentity(username)
		if true {
			log.Printf("User %s does not exist, registering new user", username)
			_, err = mspClient.Register(&msp.RegistrationRequest{
				Name:        username,
				Type:        identityTypeUser,
				Attributes:  testAttributes,
				Affiliation: orgName,
				Secret:      secret,
			})
		} else {
			log.Printf("Identity: %s", identity.Secret)
		}
		//enroll user
		err = mspClient.Enroll(username, msp.WithSecret(secret))
		if err != nil {
			log.Printf("enroll %s failed: %v", username, err)
			return "failed " + err.Error(), false
		}

		return username + " enrolled Successfully", true
	}
	log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " already enrolled", true
}

// GetArgs get [][]byte args from string array
func GetArgs(args []string) [][]byte {
	var result [][]byte
	for _, element := range args {
		result = append(result, []byte(element))
	}
	return result
}
