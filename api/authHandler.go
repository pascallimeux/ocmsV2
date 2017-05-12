package api

import (
	"errors"
	"net/http"
	"github.com/pascallimeux/ocmsV2/helpers"
)

func GetUserCredentials(r *http.Request)(helpers.UserCredentials, error){
	userCredentials := helpers.UserCredentials{}
	username, password, ok :=r.BasicAuth()
	if ok{
		log.Debug("GetUserCredentials(user:" + username + ") : calling method -")
		userCredentials.UserName = username
		userCredentials.EnrollmentSecret = password
		return userCredentials, nil
	}
	return userCredentials, errors.New("no credential in request")
}

func InitHelper (r *http.Request, helper helpers.Helper)  error {
	userCredentials, err := GetUserCredentials(r)
	if err != nil {
		return err
	}
	err = helper.Init(userCredentials)
	return err
}