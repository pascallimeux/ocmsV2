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

	netHelper.DeployCC(CHAINCODEPATH, CHAINCODEVERSION, CHAINCODEID)
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

func TestMain(m *testing.M) {
	setup()
	time.Sleep(time.Millisecond * 3000)
	code := m.Run()
	shutdown()
	os.Exit(code)
}
