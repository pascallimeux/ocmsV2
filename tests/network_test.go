package tests

import (
	"github.com/pascallimeux/ocmsV2/helpers"
	"testing"
	"fmt"
	"os"
	"github.com/cloudflare/cfssl/log"
)

var NetHelper  helpers.NetworkHelper
var ConsHelper helpers.ConsentHelper

const(
	CHAINCODEPATH    = "github.com/consentv2"
	CHAINCODEVERSION = "v0"
	CHAINCODEID      = "consentv2"
	ADMLOGIN         = "admin"
	ADMPWD           = "adminpw"
	REPO             = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures"
	STATSTOREPATH    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/enroll_user"
	CONFIGFILE       = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/config/config.yaml"
	CHANNELCONFIG    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/channel/testchannel.tx"
	CHAINID          = "testchannel"
	PROVIDERNAME     = "SW"
)

func TestDeployCC(t *testing.T) {
	err := NetHelper.DeployCC(CHAINCODEPATH, CHAINCODEVERSION, CHAINCODEID)
	if err != nil {
		fmt.Errorf("DeployCC return error: %v", err)
	}
}



func TestMain(m *testing.M) {

	NetHelper = helpers.NetworkHelper{
		Repo:                   REPO,
		ConfigFile:      	CONFIGFILE,
		ChannelConfig:   	CHANNELCONFIG,
		ChainID:         	CHAINID,
	}

	err := NetHelper.InitNetwork(ADMLOGIN, ADMPWD, STATSTOREPATH, PROVIDERNAME)
	if err != nil {
		log.Fatal(err)
	}

	ConsHelper = helpers.ConsentHelper{
		ChainID:         	CHAINID,
		Chain:			NetHelper.Chain,
		EventHub:		NetHelper.EventHub,
	}

	code := m.Run()
	os.Exit(code)
}