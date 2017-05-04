package ocms

import (
	"fmt"
	"log"
	//"os"
	"crypto/x509"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fabricCAClient "github.com/hyperledger/fabric-sdk-go/fabric-ca-client"
	//kvs "github.com/hyperledger/fabric-sdk-go/fabric-client/keyvaluestore"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/config"
	fcUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	"errors"
	"time"
	"io/ioutil"
	"encoding/pem"
	"github.com/hyperledger/fabric/bccsp"
	//"path"
)

type OcmsApp struct {
	Client          	fabricClient.Client
	CaClient        	fabricCAClient.Services
	AdminUser      	 	fabricClient.User
	ConnectEventHub 	bool
	EventHub        	events.EventHub
	ChainID         	string
	Chain 	        	fabricClient.Chain
	Initialized     	bool
	ConfigFile      	string
	ChannelConfig   	string
	Repo                    string
}

func (app *OcmsApp) InitConfig() error{
	err := config.InitConfig(app.ConfigFile)
	if err != nil {
		return err
	}
	app.Initialized = true
	return nil
}

func (app *OcmsApp) Setup() error{
	err := bccspFactory.InitFactories(&bccspFactory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &bccspFactory.SwOpts{
			HashFamily: config.GetSecurityAlgorithm(),
			SecLevel:   config.GetSecurityLevel(),
			FileKeystore: &bccspFactory.FileKeystoreOpts{
				KeyStorePath: config.GetKeyStorePath(),
			},
			Ephemeral: false,
		},
	})
	if err != nil {
		return errors.New("Failed getting ephemeral software-based BCCSP [%s]"+ err.Error())
	}

	// Get client
	client, err := fcUtil.GetClient("admin", "adminpw", app.Repo+"/enroll_user")
	if err != nil {
		return fmt.Errorf("Create client failed: %v", err)
	}
	app.Client = client

	// Get clientCa
	/*caClient, err := fabricCAClient.NewFabricCAClient()
	if err != nil {
		return errors.New("NewFabricCAClient return error: %v"+ err.Error())
	}
	app.caClient = caClient
*/
	// Get chain
	chain, err := fcUtil.GetChain(app.Client, app.ChainID)
	if err != nil {
		return fmt.Errorf("Create chain (%s) failed: %v", app.ChainID, err)
	}
	app.Chain = chain

	// Create and join channel
	if err := fcUtil.CreateAndJoinChannel(app.Client, app.Chain, app.ChannelConfig); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	// Get envenHub
	eventHub, err := getEventHub()
	if err != nil {
		return err
	}

	if app.ConnectEventHub {
		if err := eventHub.Connect(); err != nil {
			return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
		}
	}
	app.EventHub = eventHub
	app.Initialized = true
	fmt.Println("setup OK...")
	return nil
}

func (app *OcmsApp) getUser(username, password string) (fabricClient.User, error) {
	fmt.Println("---Get user %s:"+username)
	user, err := app.Client.LoadUserFromStateStore(username)

	if err != nil {
		return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
	}
	if user == nil {
		fmt.Println("---Enroll the user %s:"+username)
		key, cert, err := app.CaClient.Enroll(username, password)
		if err != nil {
			return user, errors.New("Enroll return error: %v"+ err.Error())
		}
		if key == nil {
			return user, errors.New("private key return from Enroll is nil")
		}
		if cert == nil {
			return user, errors.New("cert return from Enroll is nil")
		}

		certPem, _ := pem.Decode(cert)
		if err != nil {
			return user, errors.New("pem Decode return error: %v"+ err.Error())
		}

		cert509, err := x509.ParseCertificate(certPem.Bytes)
		if err != nil {
			return user, errors.New("x509 ParseCertificate return error: %v"+ err.Error())
		}
		if cert509.Subject.CommonName != username {
			return user, errors.New("CommonName in x509 cert is not the enrollmentID")
		}

		keyPem, _ := pem.Decode(key)
		if err != nil {
			return user, errors.New("pem Decode return error: %v"+ err.Error())
		}
		user = fabricClient.NewUser(username)
		k, err := app.Client.GetCryptoSuite().KeyImport(keyPem.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
		if err != nil {
			return user, errors.New("KeyImport return error: %v"+ err.Error())
		}
		user.SetPrivateKey(k)
		user.SetEnrollmentCertificate(cert)
		err = app.Client.SaveUserToStateStore(user, false)
		if err != nil {
			return user, errors.New("client.SetUserContext return error: %v"+ err.Error())
		}
		user, err = app.Client.LoadUserFromStateStore(username)
		if err != nil {
			return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
		}
		if user == nil {
			return user, errors.New("client.GetUserContext return nil")
		}
	}
	return user, nil
}


func (app *OcmsApp) registerUser(adminUser fabricClient.User, userName string) (string, error){
	registerRequest := fabricCAClient.RegistrationRequest{Name: userName, Type: "user", Affiliation: "org1.department1"}
	enrolmentSecret, err := app.CaClient.Register(adminUser, &registerRequest)
	if err != nil {
		return enrolmentSecret, errors.New("Error from Register: %s"+ err.Error())
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", userName, enrolmentSecret)
	return enrolmentSecret, nil
}

func (app *OcmsApp) enrollUser(userName, enrolmentSecret string) error{
	key, cert, err := app.CaClient.Enroll(userName, enrolmentSecret)
	if err != nil {
		return errors.New("Error enroling user: %s"+ err.Error())
	}
	err = ioutil.WriteFile("/tmp/"+userName+".cert.pem", cert, 0644)
	if err != nil {
		return errors.New("Error write certificate: %s"+ err.Error())
	}
	err = ioutil.WriteFile("/tmp/"+userName+".key.pem", key, 0644)
	if err != nil {
		return errors.New("Error write key: %s"+ err.Error())
	}
	return nil
}

func (app *OcmsApp) revokeUser(adminUser fabricClient.User, userName string)error{
	revokeRequest := fabricCAClient.RevocationRequest{Name: userName}
	err := app.CaClient.Revoke(adminUser, &revokeRequest)
	if err != nil {
		return errors.New("Error from Revoke: %s"+ err.Error())
	}
	return nil
}


func (app *OcmsApp) DeployCC(chainCodePath, chainCodeVersion, chainCodeID string) error {
	if err := app.InstallCC(chainCodePath, chainCodeVersion, chainCodeID, nil); err != nil {
		return err
	}
	var args []string
	return app.InstantiateCC(chainCodePath, chainCodeVersion, chainCodeID, args)
}

func (app *OcmsApp) InstallAndInstantiateExampleCC(chainCodePath, chainCodeVersion, chainCodeID string ) error {
	if err := app.InstallCC(chainCodePath, chainCodeVersion, chainCodeID,  nil); err != nil {
		return err
	}
	var args []string
	args = append(args, "init")
	args = append(args, "a")
	args = append(args, "100")
	args = append(args, "b")
	args = append(args, "200")
	return app.InstantiateCC(chainCodePath, chainCodeVersion, chainCodeID, args)
}

func (app *OcmsApp) MoveFundsExample(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "1")
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponse, txID, err := fcUtil.CreateAndSendTransactionProposal(app.Chain, chainCodeID, app.ChainID, args, []fabricClient.Peer{app.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}

	// Register for commit event
	done, fail := fcUtil.RegisterTxEvent(txID, app.EventHub)

	_, err = fcUtil.CreateAndSendTransaction(app.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}

	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}
	return txID, nil
}

func (app *OcmsApp) QueryAssetExample(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "b")
	return app.Query(chainCodeID, args)
}

func (app *OcmsApp) GetVersionExample(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "version")
	return app.Query(chainCodeID, args)
}

func (app *OcmsApp) InstantiateCC(chainCodePath, chainCodeVersion, chainCodeID string, args []string) error {
	if err := fcUtil.SendInstantiateCC(app.Chain, chainCodeID, app.ChainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{app.Chain.GetPrimaryPeer()}, app.EventHub); err != nil {
		return err
	}
	fmt.Println("Instantiate OK...")
	return nil
}

func (app *OcmsApp) InstallCC(chainCodePath, chainCodeVersion, chainCodeID string, chaincodePackage []byte) error {
	if err := fcUtil.SendInstallCC(app.Client, app.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, app.Chain.GetPeers(), app.Repo); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}
	fmt.Println("Install OK...")
	return nil
}


// GetDeployPath
/*func (app *OcmsApp) GetDeployPath() string {
	//pwd, _ := os.Getwd()
	//return path.Join(pwd, "../fixtures")
	return REPO
}*/

// getEventHub initilizes the event hub
func getEventHub() (events.EventHub, error) {
	eventHub := events.NewEventHub()
	foundEventHub := false
	peerConfig, err := config.GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading peer config: %v", err)
	}
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("******* EventHub connect to peer (%s:%d) *******\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, fmt.Errorf("No EventHub configuration found")
	}

	return eventHub, nil
}


func main() {
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//os.Setenv("OCMSPATH", dir)
	//os.Setenv("OCMSPATH", "/opt/gopath/src/github.com/pascallimeux/ocmsV2")
	//fmt.Println("OCMSPATH:", os.Getenv("OCMSPATH"))

	app := OcmsApp{
		Repo:                   "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures",
		ConfigFile:      	"/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/config/config.yaml",
		ChannelConfig:   	"/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/channel/testchannel.tx",
		ChainID:         	"testchannel",
		ConnectEventHub:        true,
	}

	err := app.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = app.Setup()
	if err != nil {
		log.Fatal(err)
	}
	/*adminUser, err := app.getUser("admin", "adminpw")
	if err != nil {
		log.Fatal(err)
	}
	username :="pascal10"
	enrolmentSecret, err := app.registerUser(adminUser, username)
	if err != nil {
		log.Fatal(err)
	}
	err = app.enrollUser(username, enrolmentSecret)
	if err != nil {
		log.Fatal(err)
	}
	err = app.revokeUser(adminUser, username)
	if err != nil {
		log.Fatal(err)
	}*/
	//_,err = app.getUser(username, enrolmentSecret)
	//if err != nil {
	//	log.Fatal(err)
	//}

	chainCodePath := "github.com/example_cc"
	chainCodeVersion := "v0"
	chainCodeID := "exemple"
	err = app.InstallAndInstantiateExampleCC(chainCodePath, chainCodeVersion, chainCodeID)
	//err = app.InstallAndInstantiateCC()
	if err != nil {
		fmt.Errorf("Install and instanciate return error: %v", err)
	}


	chaincodeQueryResponse, err := app.Client.QueryInstalledChaincodes(app.Chain.GetPrimaryPeer())
	if err != nil {
		fmt.Errorf("QueryInstalledChaincodes return error: %v", err)
	}

	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		fmt.Printf("Found deployed chaincode: %s\n", chaincode)
	}


	//time.Sleep(time.Duration(5)*time.Second)

	value, err := app.QueryAssetExample(chainCodeID)
	if err != nil {
		fmt.Errorf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)
/*
	eventID := "test([a-zA-Z]+)"

	// Register callback for chaincode event
	done, rce := fcUtil.RegisterCCEvent(app.ChainCodeID, eventID, app.EventHub)

	_, err = app.MoveFundsExample()
	if err != nil {
		fmt.Errorf("Move funds return error: %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second * 20):
		fmt.Errorf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	app.EventHub.UnregisterChaincodeEvent(rce)

	valueAfterInvoke, err := app.QueryAssetExample()
	if err != nil {
		fmt.Errorf("getQueryValue return error: %v", err)
		return
	}
	fmt.Printf("*** QueryValue after invoke %s\n", valueAfterInvoke)

	version, err := app.GetVersion()
	if err != nil {
		fmt.Errorf("Get version return error: %v", err)
	}
	fmt.Println("version: "+version)

	txID, err := app.CreateConsent()
	if err != nil {
		fmt.Errorf("Create consent return error: %v", err)
	}
	fmt.Println("txID: "+txID)

	consents, err := app.GetConsents()
	if err != nil {
		fmt.Errorf("Get consents return error: %v", err)
	}
	fmt.Println("consents: "+consents)
*/
	// voir pourquoi le SC consent declenche un PB de MSP...
	// comment fonctionne register enroll revoke
	// comment supprimer un user register
	// probleme de path pour les fichiers de config et certificats
}
const (
	APPID1      = "APP4TESTS1"
	OWNERID1    = "owner1"
	CONSUMERID1 = "consumer1"
	DATATYPE1   = "type1"
	DATAACCESS1 = "access1"
)

func (app *OcmsApp) CreateConsent(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "postconsent")
	args = append(args, APPID1)
	args = append(args, OWNERID1)
	args = append(args, CONSUMERID1)
	args = append(args, DATATYPE1)
	args = append(args, DATAACCESS1)
	args = append(args, getStringDateNow(0))
	args = append(args, getStringDateNow(7))
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponse, txID, err := fcUtil.CreateAndSendTransactionProposal(app.Chain, chainCodeID, app.ChainID, args, []fabricClient.Peer{app.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}

	// Register for commit event
	done, fail := fcUtil.RegisterTxEvent(txID, app.EventHub)

	_, err = fcUtil.CreateAndSendTransaction(app.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}

	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}
	return txID, nil
}

func (app *OcmsApp) GetConsents(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "getconsents")
	args = append(args, APPID1)
	return app.Query(chainCodeID, args)
}

func (app *OcmsApp) GetVersion(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "getversion")
	return app.Query(chainCodeID, args)
}

func (app *OcmsApp) Query(chainCodeID string, args []string) (string, error) {
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponses, _, err := fcUtil.CreateAndSendTransactionProposal(app.Chain, chainCodeID, app.ChainID, args, []fabricClient.Peer{app.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	return string(transactionProposalResponses[0].GetResponsePayload()), nil
}

func getStringDateNow(nbdaysafter time.Duration) string{
	t := time.Now().Add(nbdaysafter * 24 * time.Hour)
	return t.Format("2006-01-02")
}
