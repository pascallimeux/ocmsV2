package helpers

import (
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fabricCAClient "github.com/hyperledger/fabric-sdk-go/fabric-ca-client"
	"errors"
	"github.com/google/certificate-transparency/go/x509"
	"github.com/hyperledger/fabric/bccsp"
	"encoding/pem"
)

type UserHelper struct {
	AdmClient       fabricClient.Client
	AdmUser		fabricClient.User
	CaClient       	fabricCAClient.Services
}

type RegistrerUser struct {
	Name 		string
	Type 		string
	Affiliation     string
}

func (uh *UserHelper) init (usernameAdm, passwordAdm string) error{
	caClient, err := fabricCAClient.NewFabricCAClient()
	if err != nil {
		return errors.New("NewFabricCAClient return error: %v"+ err.Error())
	}
	uh.CaClient = caClient
	AdmUser, err := uh.getUser(usernameAdm, passwordAdm)
	if err != nil {
		return err
	}
	uh.AdmUser = AdmUser
	return nil
}

func (uh *UserHelper) registerUser(registerUser RegistrerUser) (string, error) {
	log.Debug("registerUser(name:"+ registerUser.Name+" Type:" + registerUser.Type +" Affiliation:"+ registerUser.Affiliation+") : calling method -")

	registerRequest := fabricCAClient.RegistrationRequest{Name: registerUser.Name, Type: registerUser.Type, Affiliation: registerUser.Affiliation}
	enrolmentSecret, err := uh.CaClient.Register(uh.AdmUser, &registerRequest)
	if err != nil {
		return "", err
	}
	return enrolmentSecret, nil
}

func (uh *UserHelper) getUser(username, password string) (fabricClient.User, error) {
	log.Debug("getUser(username:"+ username+") : calling method -")
	user, err := uh.AdmClient.LoadUserFromStateStore(username)

	if err != nil {
		return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
	}
	if user == nil {
		log.Debug("---Enroll the user %s:"+username)
		key, cert, err := uh.CaClient.Enroll(username, password)
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
		user, err = uh.AdmClient.LoadUserFromStateStore(username)
		if err != nil {
			return user, errors.New("client.GetUserContext return error: %v"+ err.Error())
		}
		if user == nil {
			return user, errors.New("client.GetUserContext return nil")
		}
	}
	return user, nil
}