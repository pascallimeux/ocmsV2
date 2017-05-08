package tests

import (
	"github.com/pascallimeux/ocmsV2/helpers"
	"testing"
	"os"
	"time"
	"github.com/op/go-logging"
)

var netHelper  helpers.NetworkHelper
var consHelper helpers.ConsentHelper
var log = logging.MustGetLogger("tests")


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
	CCVERSION   	 = "Orange Consent Application chaincode ver 3 Dated 2017-03-09"

	APPID1      	 = "APP4TESTS1"
	APPID2     	 = "APP4TESTS2"
	APPID3     	 = "APP4TESTS3"
	APPID4     	 = "APP4TESTS4"
	APPID5     	 = "APP4TESTS5"
	APPID6     	 = "APP4TESTS6"
	OWNERID1   	 = "owner1"
	OWNERID2   	 = "owner2"
	OWNERID3   	 = "owner3"
	CONSUMERID1	 = "consumer1"
	CONSUMERID2	 = "consumer2"
	CONSUMERID3	 = "consumer3"
	DATATYPE1  	 = "type1"
	DATAACCESS1	 = "access1"
)


func setup() {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	netHelper = helpers.NetworkHelper{
		Repo:                   REPO,
		ConfigFile:      	CONFIGFILE,
		ChannelConfig:   	CHANNELCONFIG,
		ChainID:         	CHAINID,
	}

	err := netHelper.InitNetwork(ADMLOGIN, ADMPWD, STATSTOREPATH, PROVIDERNAME)
	if err != nil {
		log.Fatal(err.Error())
	}

	consHelper = helpers.ConsentHelper{
		ChainID:         	CHAINID,
		Chain:			netHelper.Chain,
		EventHub:		netHelper.EventHub,
	}

	//netHelper.DeployCC(CHAINCODEPATH, CHAINCODEVERSION, CHAINCODEID)

	/*err = NetHelper.DeployCC(CHAINCODEPATH, CHAINCODEVERSION, CHAINCODEID)
	if err != nil {
		log.Fatal(err.Error())
	}*/
}

func shutdown(){
	_, err := consHelper.DeleteConsents4Application(CHAINCODEID, APPID1)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID2)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID3)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID4)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID5)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID6)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
}

func TestQueryInfos(t *testing.T) {
	t.Log("TestQueryInfos ")
	bci, err := netHelper.QueryInfos()
	if err != nil {
		t.Error(err)
	}
	t.Log("bcinfo string: ", bci.String())
}


func TestQueryTransaction(t *testing.T) {
	txID := createtransaction(t)
	time.Sleep(time.Millisecond * 1500)
	processedTransaction, err :=netHelper.QueryTransaction(txID)
	if err != nil {
		t.Error("QueryTransaction return error: ", err)
	}
	t.Log("transaction: ", processedTransaction.String())
}

func TestQueryBlockByNumber(t *testing.T) {
	block, err := netHelper.QueryBlockByNumber(1)
	if err != nil {
		t.Error("QueryBlockByNumber return error: ", err)
	}
	t.Log("block: ", block.String())
}

func TestQueryBlockByHash(t *testing.T) {
	bci, err := netHelper.QueryInfos()
	if err != nil {
		t.Fatalf("QueryInfo return error: %v", err)
	}

	// Test Query Block by Hash - retrieve current block by hash
	block, err := netHelper.QueryBlockByHash(bci.CurrentBlockHash)
	if err != nil {
		t.Fatalf("QueryBlockByHash return error: %v", err)
	}
	t.Log("block: ", block.String())
}

func TestQueryChannels(t *testing.T) {
	channelQueryResponse, err := netHelper.QueryChannels()
	if err != nil {
		t.Fatalf("QueryChannels return error: %v", err)
	}
	for _, channel := range channelQueryResponse.Channels {
		t.Log("Channel: ",channel, "\n")
	}
}

func TestGetInstalledChainCode(t *testing.T) {
	chaincodeQueryResponse, err := netHelper.GetInstalledChainCode()
	if err != nil {
		t.Fatalf("QueryInstalledChaincodes return error: %v", err)
	}

	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		t.Log("InstalledCC: ",chaincode,"\n")
	}
}

func TestQueryByChainCode(t *testing.T) {
	chaincodeQueryResponse, err := netHelper.QueryByChainCode("lccc")
	if err != nil {
		t.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}

	for _, chaincode := range chaincodeQueryResponse {
		t.Log("InstantiatedCC: ",chaincode,"\n")
	}
}



func TestMain(m *testing.M) {
	setup()
	time.Sleep(time.Millisecond * 3000)
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func createtransaction(t *testing.T)string{
	txID, err := consHelper.CreateConsent(CHAINCODEID, APPID1, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	return txID
}