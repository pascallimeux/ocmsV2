package helpers

import (
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
)

type UserHelper struct {
	ChainID string
	Chain 	fabricClient.Chain
	Repo    string
}
