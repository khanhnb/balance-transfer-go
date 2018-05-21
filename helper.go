package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

func getRegisteredUser(username, orgName string) string {
	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	testAttributes := []msp.Attribute{
		{
			Name:  "1",
			Value: fmt.Sprintf("%s:ecert", "2"),
			ECert: true,
		},
		{
			Name:  "3",
			Value: fmt.Sprintf("%s:ecert", "4"),
			ECert: true,
		},
	}

	registrarEnrollID, registrarEnrollSecret := getRegistrarEnrollmentCredentials(ctxProvider)

	err = mspClient.Enroll(registrarEnrollID, msp.WithSecret(registrarEnrollSecret))
	if err != nil {
		log.Fatalf("enroll registrar failed: %v", err)
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
	if err != nil {
		log.Fatalf(":::::: %v\n", err)
	}

	// Enroll the new user
	err = mspClient.Enroll(username, msp.WithSecret(enrollmentSecret))
	log.Printf("secret: %s", enrollmentSecret)
	if err != nil {
		log.Printf("enroll %s failed: %v", username, err)
	}
	return ""
}

func getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {

	ctx, err := ctxProvider()
	if err != nil {
		fmt.Printf("failed to get context: %v\n", err)
	}

	clientConfig, err := ctx.IdentityConfig().Client()
	if err != nil {
		fmt.Printf("config.Client() failed: %v\n", err)
	}

	myOrg := clientConfig.Organization

	caConfig, err := ctx.IdentityConfig().CAConfig(myOrg)
	if err != nil {
		fmt.Printf("CAConfig failed: %v\n", err)
	}

	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}
