package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	var config = config.FromFile("config.yaml")
	fabsdk.New(config)
}
