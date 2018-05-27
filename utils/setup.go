package utils

import (
	"fmt"
	"log"
	"os"

	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
)

// FabricSetup struct
type FabricSetup struct {
	AdminUser         string
	OrdererOrgName    string
	ConfigFileName    string
	Secret            []byte
	IdentityTypeUser  string
	Sdk               *fabsdk.FabricSDK
	RegistrarUsername string
	RegistrarPassword string
	MspClient         *clientMSP.Client
}

// Init reads config file, setup client, CA
func (hfc *FabricSetup) Init() {
	hfc.Sdk = hfc.setupSDK(hfc.ConfigFileName)
	// hfc.cleanupUserData()
	hfc.setupCA()
}

func (hfc *FabricSetup) setupSDK(configFileName string) *fabsdk.FabricSDK {
	var config = config.FromFile(configFileName)
	sdk, err := fabsdk.New(config)
	if err != nil {
		fmt.Printf("failed to create new SDK: %s\n", err)
	}
	return sdk
}

func (hfc *FabricSetup) setupCA() {
	ctxProvider := hfc.Sdk.Context()
	mspClient, err := clientMSP.New(ctxProvider)
	registrarEnrollID, registrarEnrollSecret := hfc.getRegistrarEnrollmentCredentials(ctxProvider)

	err = mspClient.Enroll(registrarEnrollID, clientMSP.WithSecret(registrarEnrollSecret))
	if err != nil {
		log.Fatalf("enroll registrar failed: %v", err)
	}
	hfc.MspClient = mspClient
}

func (hfc *FabricSetup) cleanupUserData() {
	configBackend, err := hfc.Sdk.Config()
	if err != nil {
		log.Fatal(err)
	}

	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := msp.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatal(err)
	}

	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	hfc.cleanupPath(keyStorePath)
	hfc.cleanupPath(credentialStorePath)
}

func (hfc *FabricSetup) cleanupPath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		log.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}

func (hfc *FabricSetup) getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {

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
