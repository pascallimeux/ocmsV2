package helpers

import (
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fabricCAClient "github.com/hyperledger/fabric-sdk-go/fabric-ca-client"
	"errors"
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	"github.com/google/certificate-transparency/go/x509"
	"github.com/hyperledger/fabric/bccsp"
	"encoding/pem"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

type UserHelper struct {
	StatStorePath   string
	AdmClient       fabricClient.Client
	AdmUser		fabricClient.User
	CaClient       	fabricCAClient.Services
	Initialized	bool
}


type UserCredentials struct {
	UserName 	 string	   `json:"username"`
	EnrollmentSecret string	   `json:"password"`
}

type UserRegistrer struct {
	Name 		 string	`json:"name"`
	Type 		 string	`json:"type"`
	Affiliation      string	`json:"affiliation"`
}

func (uh *UserHelper) Init (userCredentials UserCredentials) error{
	client, err := sdkUtil.GetClient("admin", "adminpw", uh.StatStorePath)
	if err != nil {
		return errors.New("Create client failed: %v" + err.Error())
	}
	uh.AdmClient = client
	caClient, err := fabricCAClient.NewFabricCAClient()
	if err != nil {
		return errors.New("NewFabricCAClient return error: %v" + err.Error())
	}
	uh.CaClient = caClient
	AdmUser, err := uh.GetUser(userCredentials)
	if err != nil {
		return err
	}
	uh.AdmUser = AdmUser
	uh.Initialized = true
	return nil
}

func (uh *UserHelper) RegisterUser(registerUser UserRegistrer) (string, error) {
	log.Debug("registerUser(name:"+ registerUser.Name+" Type:" + registerUser.Type +" Affiliation:"+ registerUser.Affiliation+") : calling method -")
	registerRequest := fabricCAClient.RegistrationRequest{Name: registerUser.Name, Type: registerUser.Type, Affiliation: registerUser.Affiliation}
	enrolmentSecret, err := uh.CaClient.Register(uh.AdmUser, &registerRequest)
	if err != nil {
		return "", err
	}
	return enrolmentSecret, nil
}


func (uh *UserHelper) EnrollUser(userCredentials UserCredentials) error{
	log.Debug("enrollUser(userName:"+ userCredentials.UserName+") : calling method -")
	key, cert, err := uh.CaClient.Enroll(userCredentials.UserName, userCredentials.EnrollmentSecret)
	if err != nil {
		return errors.New("Error enroling user: %s"+ err.Error())
	}
	err = ioutil.WriteFile(uh.StatStorePath+"/"+userCredentials.UserName+".cert.pem", cert, 0644)
	if err != nil {
		return errors.New("Error write certificate: %s"+ err.Error())
	}
	err = ioutil.WriteFile(uh.StatStorePath+"/"+userCredentials.UserName+".key.pem", key, 0644)
	if err != nil {
		return errors.New("Error write key: %s"+ err.Error())
	}
	return nil
}

func (uh *UserHelper) ReenrollUser(userCredentials UserCredentials) error{
	log.Debug("ReenrollUser(userName:"+ userCredentials.UserName+") : calling method -")
	enrolleduser := fabricClient.NewUser(userCredentials.UserName)
	//enrolleduser.SetEnrollmentCertificate(ecert)
	//enrolleduser.SetPrivateKey(k)
	key, cert, err := uh.CaClient.Reenroll(enrolleduser)
	if err != nil {
		return errors.New("Error Reenroling user: %s" + err.Error())
	}
	err = ioutil.WriteFile(uh.StatStorePath+userCredentials.UserName+".cert.pem", cert, 0644)
	if err != nil {
		return errors.New("Error write certificate: %s"+ err.Error())
	}
	err = ioutil.WriteFile(uh.StatStorePath+userCredentials.UserName+".key.pem", key, 0644)
	if err != nil {
		return errors.New("Error write key: %s"+ err.Error())
	}
	return nil
}

func (uh *UserHelper) RevokeUser(userCredentials UserCredentials)error{
	log.Debug("revokeUser(userName:"+ userCredentials.UserName+") : calling method -")
	revokeRequest := fabricCAClient.RevocationRequest{Name: userCredentials.UserName}
	err := uh.CaClient.Revoke(uh.AdmUser, &revokeRequest)
	if err != nil {
		return errors.New("Error from Revoke: %s"+ err.Error())
	}
	return nil
}

func (uh *UserHelper) GetUser(userCredentials UserCredentials) (fabricClient.User, error) {
	log.Debug("getUser(username:"+ userCredentials.UserName+") : calling method -")
	user, err := uh.AdmClient.LoadUserFromStateStore(userCredentials.UserName)
	if err != nil {
		return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
	}
	if user == nil {
		log.Debug("---Enroll the user %s:"+userCredentials.UserName)
		key, cert, err := uh.CaClient.Enroll(userCredentials.UserName, userCredentials.EnrollmentSecret)
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
		if cert509.Subject.CommonName != userCredentials.UserName {
			return user, errors.New("CommonName in x509 cert is not the enrollmentID")
		}

		keyPem, _ := pem.Decode(key)
		if err != nil {
			return user, errors.New("pem Decode return error: %v"+ err.Error())
		}
		user = fabricClient.NewUser(userCredentials.UserName)
		k, err := uh.AdmClient.GetCryptoSuite().KeyImport(keyPem.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
		if err != nil {
			return user, errors.New("KeyImport return error: %v"+ err.Error())
		}
		user.SetPrivateKey(k)
		user.SetEnrollmentCertificate(cert)
		err = uh.AdmClient.SaveUserToStateStore(user, false)
		if err != nil {
			return user, errors.New("client.SetUserContext return error: %v"+ err.Error())
		}
		user, err = uh.AdmClient.LoadUserFromStateStore(userCredentials.UserName)
		if err != nil {
			return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
		}
		if user == nil {
			return user, errors.New("client.GetUserContext return nil")
		}
	}
	return user, nil
}


func getClient(userCredentials UserCredentials, statStorePath string) (fabricClient.Client, error) {
	log.Debug("GetClient(username:"+ userCredentials.UserName+") : calling method -")
	client, err := sdkUtil.GetClient(userCredentials.UserName, userCredentials.EnrollmentSecret, statStorePath)
	if err != nil {
		log.Debug("getClient return error: %v" + err.Error())
		return client, errors.New("getClient return error: %v" + err.Error())
	}
	return client, nil
}

func CreateRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return "user" + strconv.Itoa(rand.Intn(500000))
}
