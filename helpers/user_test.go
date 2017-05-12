package helpers

import (
	"testing"
)

func TestRegisterUser(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	if enrollSecret == ""{
		t.Error("no enrollSecret received")
	}
}

func TestEnrollUser(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	userCredentials := UserCredentials{UserName: username, EnrollmentSecret: enrollSecret}
	err = userHelper.EnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
}

func TestRevokeUser(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	userCredentials := UserCredentials{UserName: username, EnrollmentSecret: enrollSecret}
	err = userHelper.EnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
	err = userHelper.RevokeUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	userCredentials := UserCredentials{UserName: username, EnrollmentSecret: enrollSecret}
	user, err := userHelper.GetUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
	t.Log(user.GetName())
}


func TestGetClient(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	userCredentials := UserCredentials{UserName: username, EnrollmentSecret: enrollSecret}
	_, err = getClient(userCredentials, statStorePath)
	if err != nil {
		t.Error(err)
	}
}

func __TestReenrollUser(t *testing.T) {
	username := CreateRandomName()
	registerUser := UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1"}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	t.Log("enrollSecret string: ",enrollSecret)
	userCredentials := UserCredentials{UserName: username, EnrollmentSecret: enrollSecret}
	err = userHelper.EnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
	err = userHelper.RevokeUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
	err = userHelper.EnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
}
