package main

import(
	"github.com/pascallimeux/ocmsV2/helpers"
	"github.com/cloudflare/cfssl/log"
	"fmt"
)


const(
	chainCodePath    = "github.com/consentv2"
	chainCodeVersion = "v0"
	chainCodeID      = "consentv2"
	adminusername    = "admin"
	adminPwd         = "adminpw"
	repo             = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures"
	statstorePath    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/enroll_user"
	configfile       = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/config/config.yaml"
	channelConfig    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/channel/testchannel.tx"
	chainID          = "testchannel"
	providerName     = "SW"
)

func main() {

	networkHelper := helpers.NetworkHelper{
		Repo:                   repo,
		ConfigFile:      	configfile,
		ChannelConfig:   	channelConfig,
		ChainID:         	chainID,
	}

	err := networkHelper.InitNetwork(adminusername, adminPwd, statstorePath, providerName)
	if err != nil {
		log.Fatal(err)
	}

	err = networkHelper.DeployCC(chainCodePath, chainCodeVersion, chainCodeID)
	if err != nil {
		log.Fatal(err)
	}

	consentHelper := helpers.ConsentHelper{
		ChainID:         	chainID,
		Chain:			networkHelper.Chain,
		EventHub:		networkHelper.EventHub,
	}

	version, err := consentHelper.GetVersion(chainCodeID)
	if err != nil {
		fmt.Errorf("getVersion return error: %v", err)
	}
	fmt.Printf("*** Version of consentV2 CC: %s\n", version)


}
