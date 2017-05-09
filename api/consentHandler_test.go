package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocmsV2/helpers"
	"github.com/pascallimeux/ocmsV2/settings"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"strings"
	"io/ioutil"
)

var configuration settings.Settings
var httpServerTest *httptest.Server
const TransactionTimeout = time.Millisecond * 1500
const APPID = "apptest"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// Init settings
	var err error
	configuration, err = settings.GetSettings("..", "ocms")
	if err != nil {
		panic(err.Error())
	}

	netHelper := helpers.NetworkHelper{
		Repo:                   configuration.Repo,
		ConfigFile:      	configuration.SDKConfigfile,
		ChannelConfig:   	configuration.ChannelConfigFile,
		ChainID:         	configuration.ChainID,
	}

	err = netHelper.InitNetwork(configuration.Adminusername, configuration.AdminPwd, configuration.StatstorePath, configuration.ProviderName)
	if err != nil {
		log.Fatal(err.Error())
	}

	consHelper := helpers.ConsentHelper{
		ChainID:         	configuration.ChainID,
		Chain:			netHelper.Chain,
		EventHub:		netHelper.EventHub,
	}

	netHelper.DeployCC(configuration.ChainCodePath, configuration.ChainCodeVersion, configuration.ChainCodeID)
	/*err = netHelper.DeployCC(configuration.ChainCodePath, configuration.ChainCodeVersion, configuration.ChainCodeID)
	if err != nil {
		log.Fatal(err.Error())
	}*/

	// Init application context
	appContext := AppContext{ConsentHelper: consHelper, NetworkHelper: netHelper, ChainCodeID : configuration.ChainCodeID}

	router := mux.NewRouter().StrictSlash(false)
	appContext.CreateOCMSRoutes(router)

	// Init routes for application
	appContext.CreateOCMSRoutes(router)

	// Init http server for tests
	httpServerTest = httptest.NewServer(router)

}

func shutdown() {
	defer httpServerTest.Close()
	defer configuration.Close()
}

func TestCreateConsentFromAPINominal(t *testing.T) {
	consent := helpers.Consent{OwnerID: "1111", ConsumerID: "2222"}
	consentID, err := createConsent(consent)
	if err != nil {
		t.Error(err)
	}
	if consentID == "" {
		t.Error("bad consent ID")
	}
}

func TestGetConsentDetailFromAPINominal(t *testing.T) {
	consent := helpers.Consent{OwnerID: "OOOO", ConsumerID: "AAAA"}
	consentID, err := createConsent(consent)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consent2, err2 := getConsent(consentID)
	if err2 != nil {
		t.Error(err2)
	}
	if consent2.ConsentID != consentID || consent2.ConsumerID != consent.ConsumerID {
		t.Error(err)
	}

}

func TestGetConsentsFromAPINominal(t *testing.T) {
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "2222"})
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "3333"})
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "4444"})
	consents, err := getListOfConsents("", "")
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetConsents4OwnerFromAPINominal(t *testing.T) {
	ownerid := "1111"
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "2222"})
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "3333"})
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: "4444"})
	consents, err := getListOfConsents(ownerid, "")
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetConsents4ConsumerFromAPINominal(t *testing.T) {
	consumerid := "3333"
	createConsent(helpers.Consent{OwnerID: "1111", ConsumerID: consumerid})
	createConsent(helpers.Consent{OwnerID: "2222", ConsumerID: consumerid})
	createConsent(helpers.Consent{OwnerID: "3333", ConsumerID: consumerid})
	consents, err := getListOfConsents("", consumerid)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}



func createConsent(consent helpers.Consent) (string, error) {
	var responseConsent helpers.Consent
	consent.Action = "create"
	consent.AppID = APPID
	data, _ := json.Marshal(consent)
	request, err1 := buildRequest("POST", httpServerTest.URL+CONSENTAPI, string(data))
	if err1 != nil {
		return "", err1
	}
	status, body_bytes, err2 := executeRequest(request)
	if err2 != nil {
		return "", err2
	}
	err3 := json.Unmarshal(body_bytes, &responseConsent)
	if err3 != nil {
		return "", err3
	}

	if status != http.StatusOK {
		return "", errors.New("bad status")
	}
	return responseConsent.ConsentID, nil
}

func getConsent(consentID string) (helpers.Consent, error) {
	consent := helpers.Consent{Action: "get", AppID: APPID, ConsentID: consentID}
	responseConsent := helpers.Consent{}
	data, _ := json.Marshal(consent)
	request, err1 := buildRequest("POST", httpServerTest.URL+CONSENTAPI, string(data))
	if err1 != nil {
		return responseConsent, err1
	}
	status, body_bytes, err2 := executeRequest(request)
	if err2 != nil {
		return responseConsent, err2
	}
	err3 := json.Unmarshal(body_bytes, &responseConsent)
	if err3 != nil {
		return responseConsent, err3
	}
	if status != http.StatusOK {
		return responseConsent, errors.New("bad status")
	}
	return responseConsent, nil
}

func getListOfConsents(ownerID, consumerID string) ([]helpers.Consent, error) {
	consent := helpers.Consent{Action: "list", AppID: APPID}
	consents := []helpers.Consent{}
	if ownerID != "" {
		consent.OwnerID = ownerID
		consent.Action = "list4owner"
	} else if consumerID != "" {
		consent.ConsumerID = consumerID
		consent.Action = "list4consumer"
	}
	data, _ := json.Marshal(consent)
	request, err1 := buildRequest("POST", httpServerTest.URL+CONSENTAPI, string(data))
	if err1 != nil {
		return consents, err1
	}
	status, body_bytes, err2 := executeRequest(request)
	if err2 != nil {
		return consents, err2
	}
	err3 := json.Unmarshal(body_bytes, &consents)
	if err3 != nil {
		return consents, err3
	}
	if status != http.StatusOK {
		return consents, errors.New("bad status")
	}
	return consents, nil
}

func buildRequest(method, uri, data string) (*http.Request, error) {
	var requestData *strings.Reader
	if data != "" {
		requestData = strings.NewReader(data)
	} else {
		requestData = nil
	}
	request, err := http.NewRequest(method, uri, requestData)
	if err != nil {
		return request, err
	}
	return request, nil
}


func executeRequest(request *http.Request) (int, []byte, error) {
	status := 0
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return status, nil, err
	}
	status = response.StatusCode
	body_bytes, err2 := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err2 != nil {
		return status, body_bytes, err2
	}
	return status, body_bytes, nil
}
