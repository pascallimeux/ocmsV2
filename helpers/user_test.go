package helpers

import (
	"testing"
)

func TestRegisterUser(t *testing.T) {
	registerUser := RegistrerUser{
		Name:           "pascal",
		Type:      	"user",
		Affiliation:   	"org1.department1",
	}
	enrollSecret, err :=userHelper.RegisterUser(registerUser)
	if err != nil {
		t.Error(err)
	}
	t.Log("enrollSecret string: ",enrollSecret)
}

func TestEnrollUser(t *testing.T) {

}

func TestRevokeUser(t *testing.T) {

}

