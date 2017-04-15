package main

import (
	"fmt"
	"log"
	"os"
	"crypto/x509"
	"path/filepath"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fabricCAClient "github.com/hyperledger/fabric-sdk-go/fabric-ca-client"
	kvs "github.com/hyperledger/fabric-sdk-go/fabric-client/keyvaluestore"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/config"
	fcUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	"errors"
	"io/ioutil"
	"encoding/pem"
	"github.com/hyperledger/fabric/bccsp"
	"path"
)

type OcmsApp struct {
	client          	fabricClient.Client
	caClient        	fabricCAClient.Services
	chainCodeID		string
	chainCodePath		string
	chainCodeVersion	string
	adminUser      	 	fabricClient.User
	connectEventHub 	bool
	eventHub        	events.EventHub
	chainID         	string
	chain 	        	fabricClient.Chain
	initialized     	bool
	configFile      	string
	channelConfig   	string
}

func (app *OcmsApp) initConfig() error{
	err := config.InitConfig(app.configFile)
	if err != nil {
		return err
	}
	app.initialized = true
	return nil
}

func (app *OcmsApp) setup() error{
	// Get Client
	client := fabricClient.NewClient()
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
	cryptoSuite := bccspFactory.GetDefault()
	client.SetCryptoSuite(cryptoSuite)
	stateStore, err := kvs.CreateNewFileKeyValueStore(os.Getenv("OCMSPATH")+"/fixtures/enroll_user")
	if err != nil {
		return errors.New("CreateNewFileKeyValueStore return error[%s]"+ err.Error())
	}
	client.SetStateStore(stateStore)
	app.client = client

	// Get clientCa
	caClient, err := fabricCAClient.NewFabricCAClient()
	if err != nil {
		return errors.New("NewFabricCAClient return error: %v"+ err.Error())
	}
	app.caClient = caClient

	// Get chain
	chain, err := fcUtil.GetChain(app.client, app.chainID)
	if err != nil {
		return fmt.Errorf("Create chain (%s) failed: %v", app.chainID, err)
	}
	app.chain = chain

	// Create and join channel
	if err := fcUtil.CreateAndJoinChannel(app.client, app.chain, app.channelConfig); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	eventHub, err := getEventHub()
	if err != nil {
		return err
	}

	if app.connectEventHub {
		if err := eventHub.Connect(); err != nil {
			return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
		}
	}
	app.eventHub = eventHub
	app.initialized = true
	return nil
}

func (app *OcmsApp) getUser(username, password string) (fabricClient.User, error) {
	fmt.Println("---Get user %s:"+username)
	user, err := app.client.GetUserContext(username)

	if err != nil {
		return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
	}
	if user == nil {
		fmt.Println("---Enroll the user %s:"+username)
		key, cert, err := app.caClient.Enroll(username, password)
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
		user = fabricClient.NewUser("admin")
		k, err := app.client.GetCryptoSuite().KeyImport(keyPem.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
		if err != nil {
			return user, errors.New("KeyImport return error: %v"+ err.Error())
		}
		user.SetPrivateKey(k)
		user.SetEnrollmentCertificate(cert)
		err = app.client.SetUserContext(user, false)
		if err != nil {
			return user, errors.New("client.SetUserContext return error: %v"+ err.Error())
		}
		user, err = app.client.GetUserContext("admin")
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
	enrolmentSecret, err := app.caClient.Register(adminUser, &registerRequest)
	if err != nil {
		return enrolmentSecret, errors.New("Error from Register: %s"+ err.Error())
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", userName, enrolmentSecret)
	return enrolmentSecret, nil
}

func (app *OcmsApp) enrollUser(userName, enrolmentSecret string) error{
	key, cert, err := app.caClient.Enroll(userName, enrolmentSecret)
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
	err := app.caClient.Revoke(adminUser, &revokeRequest)
	if err != nil {
		return errors.New("Error from Revoke: %s"+ err.Error())
	}
	return nil
}


func (app *OcmsApp) InstallAndInstantiateCC() error {
	if app.chainCodeID == "" {
		app.chainCodeID = fcUtil.GenerateRandomID()
	}
	if err := app.InstallCC(nil); err != nil {
		return err
	}
	var args []string
	return app.InstantiateCC(args)
}

func (app *OcmsApp) InstantiateCC(args []string) error {
	if err := fcUtil.SendInstantiateCC(app.chain, app.chainCodeID, app.chainID, args, app.chainCodePath, app.chainCodeVersion, []fabricClient.Peer{app.chain.GetPrimaryPeer()}, app.eventHub); err != nil {
		return err
	}
	return nil
}

func (app *OcmsApp) InstallCC(chaincodePackage []byte) error {
	if err := fcUtil.SendInstallCC(app.chain, app.chainCodeID, app.chainCodePath, app.chainCodeVersion, chaincodePackage, app.chain.GetPeers(), app.GetDeployPath()); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}
	return nil
}


// GetDeployPath
func (app *OcmsApp) GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, os.Getenv("OCMSPATH")+"/fixtures")
}

// getEventHub initilizes the event hub
func getEventHub() (events.EventHub, error) {
	eventHub := events.NewEventHub()
	foundEventHub := false
	for _, p := range config.GetPeersConfig() {
		if p.EventHost != "" && p.EventPort != "" {
			fmt.Printf("******* EventHub connect to peer (%s:%s) *******\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%s", p.EventHost, p.EventPort), p.TLSCertificate, p.TLSServerHostOverride)
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("OCMSPATH", dir)
	os.Setenv("OCMSPATH", "/home/blockchain/src/github.com/pascallimeux/ocmsv2")
	fmt.Println("OCMSPATH:", os.Getenv("OCMSPATH"))

	app := OcmsApp{
		configFile:      	os.Getenv("OCMSPATH")+"/fixtures/config/config.yaml",
		channelConfig:   	os.Getenv("OCMSPATH")+"/fixtures/channel/testchannel.tx",
		chainID:         	"testchannel",
		chainCodeVersion:	"v0",
		chainCodePath:		"github.com/consentv2",
		connectEventHub: true,
	}

	err = app.initConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = app.setup()
	if err != nil {
		log.Fatal(err)
	}
	/*adminUser, err := app.getUser("admin", "adminpw")
	if err != nil {
		log.Fatal(err)
	}
	username :="pascal14"
	enrolmentSecret, err := app.registerUser(adminUser, username)
	if err != nil {
		log.Fatal(err)
	}
	err = app.enrollUser(username, enrolmentSecret)
	if err != nil {
		log.Fatal(err)
	}*/
	//_,err = app.getUser(username, enrolmentSecret)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = app.revokeUser(adminUser, username)
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = app.InstallAndInstantiateCC()
	if err != nil {
		log.Fatal(err)
	}

}
