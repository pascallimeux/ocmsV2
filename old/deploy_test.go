package tests

import (
	"github.com/pascallimeux/ocmsV2/ocms"
	"testing"
	"fmt"
	"os"
)

var app ocms.OcmsApp
const (
	chainCodePath    = "github.com/consentv2"
	chainCodeVersion = "v0"
	chainCodeID      = "consentv2"
)


func TestChainCodeDeploy(t *testing.T) {
	err := app.DeployCC(chainCodePath, chainCodeVersion, chainCodeID)
	if err != nil {
		fmt.Errorf("Install and instanciate return error: %v", err)
	}

	err = app.InstallAndInstantiateExampleCC("github.com/example_cc", "v0", "exemple")
	if err != nil {
		fmt.Errorf("Install and instanciate return error: %v", err)
	}
	value, err := app.QueryAssetExample("exemple")
	if err != nil {
		fmt.Errorf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)
/*

	chaincodeQueryResponse, err := app.Client.QueryInstalledChaincodes(app.Chain.GetPrimaryPeer())
	if err != nil {
		fmt.Errorf("QueryInstalledChaincodes return error: %v", err)
	}

	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		fmt.Printf("Found deployed chaincode: %s\n", chaincode)
	}*/
	//time.Sleep(time.Duration(5)*time.Second)
}

func TestQueryExample(t *testing.T) {
	value, err := app.QueryAssetExample("exemple")
	if err != nil {
		fmt.Errorf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)
}

func TestGetVersionExample(t *testing.T) {
	value, err := app.GetVersionExample("exemple")
	if err != nil {
		fmt.Errorf("GetVersionExample return error: %v", err)
	}
	fmt.Println(value)
	fmt.Printf("*** version %s\n", value)
}

func TestGetConsentVersion(t *testing.T) {
	value, err := app.GetVersion(chainCodeID)
	if err != nil {
		fmt.Errorf("getVersion return error: %v", err)
	}
	fmt.Printf("*** Version of consentV2 CC: %s\n", value)
}

func TestCreateConsent(t *testing.T) {
	txid, err := app.CreateConsent(chainCodeID)
	if err != nil {
		fmt.Errorf("CreateConsent return error: %v", err)
	}
	fmt.Printf("*** Create consents: %s\n", txid)
}

func TestGetConsents(t *testing.T) {
	consents, err := app.GetConsents(chainCodeID)
	if err != nil {
		fmt.Errorf("GetConsents return error: %v", err)
	}
	fmt.Printf("*** GetConsents of consentV2 CC: %s\n", consents)
}


func TestMain(m *testing.M) {
	gopath := os.Getenv("GOPATH")
	app = ocms.OcmsApp{
		Repo:                   gopath+"/src/github.com/pascallimeux/ocmsV2/fixtures",
		ConfigFile:      	gopath+"/src/github.com/pascallimeux/ocmsV2/fixtures/config/config.yaml",
		ChannelConfig:   	gopath+"/src/github.com/pascallimeux/ocmsV2/fixtures/channel/testchannel.tx",
		ChainID:         	"testchannel",
		ConnectEventHub:        true,
	}

	err := app.InitConfig()
	if err != nil {
		fmt.Printf("error from Init %v", err)
		os.Exit(-1)
	}

	err = app.Setup()
	if err != nil {
		fmt.Printf("error from Setup %v", err)
		os.Exit(-1)
	}

	code := m.Run()
	os.Exit(code)
}