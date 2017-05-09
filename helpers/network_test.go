package helpers

import (
	"testing"
	"os"
	"time"
	"github.com/pascallimeux/ocmsV2/settings"
)

var netHelper  NetworkHelper
var consHelper ConsentHelper
var configuration settings.Settings

const(
	CCVERSION   	   = "Orange Consent Application chaincode ver 3 Dated 2017-03-09"
	APPID1      	   = "APP4TESTS1"
	APPID2     	   = "APP4TESTS2"
	APPID3     	   = "APP4TESTS3"
	APPID4     	   = "APP4TESTS4"
	APPID5     	   = "APP4TESTS5"
	APPID6     	   = "APP4TESTS6"
	OWNERID1   	   = "owner1"
	OWNERID2   	   = "owner2"
	OWNERID3   	   = "owner3"
	CONSUMERID1	   = "consumer1"
	CONSUMERID2	   = "consumer2"
	CONSUMERID3	   = "consumer3"
	DATATYPE1  	   = "type1"
	DATAACCESS1	   = "access1"
	TransactionTimeout = time.Millisecond * 1500
)


func setup() {

	var err error
	// Init settings
	configuration, err = settings.GetSettings("..", "ocmstest")
	if err != nil {
		panic(err.Error())
	}

	netHelper = NetworkHelper{
		Repo:                   configuration.Repo,
		ConfigFile:      	configuration.SDKConfigfile,
		ChannelConfig:   	configuration.ChannelConfigFile,
		ChainID:         	configuration.ChainID,
	}

	err = netHelper.InitNetwork(configuration.Adminusername, configuration.AdminPwd, configuration.StatstorePath, configuration.ProviderName)
	if err != nil {
		log.Fatal(err.Error())
	}

	consHelper = ConsentHelper{
		ChainID:         	configuration.ChainID,
		Chain:			netHelper.Chain,
		EventHub:		netHelper.EventHub,
	}

	netHelper.DeployCC(configuration.ChainCodePath, configuration.ChainCodeVersion, configuration.ChainCodeID)
	/*err = netHelper.DeployCC(configuration.ChainCodePath, configuration.ChainCodeVersion, configuration.ChainCodeID)
	if err != nil {
		log.Fatal(err.Error())
	}*/
}

func shutdown(){
	_, err := consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID1)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID2)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID3)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID4)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID5)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID6)
	if err != nil {
		log.Fatal("DeleteConsent return error: ", err)
	}
	defer configuration.Close()
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
	txID, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID1, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	return txID
}