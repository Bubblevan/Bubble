package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	contract := new(SmartContract)
	cc, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
