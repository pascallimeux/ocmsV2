package api

import (
	"github.com/pascallimeux/ocmsV2/helpers"
	"testing"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

func TestRegisterUserAPINominal(t *testing.T) {
	username := helpers.CreateRandomName()
	registerUser := helpers.UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1" }
	enrollmentSecret, err := sendRegister(registerUser)
	if err != nil {
		t.Error(err)
	}
	if enrollmentSecret == "" {
		t.Error("bad enrollmentSecret")
	}
}


func TestEnrollUserAPINominal(t *testing.T) {
	username := helpers.CreateRandomName()
	registerUser := helpers.UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1" }
	enrollmentSecret, err := sendRegister(registerUser)
	if err != nil {
		t.Error(err)
	}
	if enrollmentSecret == "" {
		t.Error("bad enrollmentSecret")
	}
	userCredentials := helpers.UserCredentials{UserName: username, EnrollmentSecret: enrollmentSecret}
	err = sendEnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
}


func TestRevokeUserAPINominal(t *testing.T) {
	username := helpers.CreateRandomName()
	registerUser := helpers.UserRegistrer{Name: username, Type: "user", Affiliation: "org1.department1" }
	enrollmentSecret, err := sendRegister(registerUser)
	if err != nil {
		t.Error(err)
	}
	if enrollmentSecret == "" {
		t.Error("bad enrollmentSecret")
	}
	userCredentials := helpers.UserCredentials{UserName: username, EnrollmentSecret: enrollmentSecret}
	err = sendEnrollUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
	err = sendRevokeUser(userCredentials)
	if err != nil {
		t.Error(err)
	}
}

func sendRevokeUser(userCredentials helpers.UserCredentials) error {
	data, _ := json.Marshal(userCredentials)
	request, err := buildRequestWithLoginPassword("POST", httpServerTest.URL+REVOKE, string(data), ADMINNAME, ADMINPWD)
	if err != nil {
		return err
	}
	status, _, err := executeRequest(request)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.New("bad status")
	}
	return nil
}

func sendEnrollUser(userCredentials helpers.UserCredentials) error {
	data, _ := json.Marshal(userCredentials)
	request, err := buildRequestWithLoginPassword("POST", httpServerTest.URL+ENROLL, string(data), ADMINNAME, ADMINPWD)
	if err != nil {
		return err
	}
	status, _, err := executeRequest(request)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.New("bad status")
	}
	return nil
}

func sendRegister(registerUser helpers.UserRegistrer) (string, error) {
	var response EnrollmentSecret
	data, _ := json.Marshal(registerUser)
	request, err1 := buildRequestWithLoginPassword("POST", httpServerTest.URL+REGISTER, string(data), ADMINNAME, ADMINPWD)
	if err1 != nil {
		return "", err1
	}
	status, body_bytes, err2 := executeRequest(request)
	if err2 != nil {
		return "", err2
	}
	err3 := json.Unmarshal(body_bytes, &response)
	if err3 != nil {
		return "", err3
	}
	if status != http.StatusOK {
		return "", errors.New("bad status")
	}
	return response.Secret, nil
}

func buildRequestWithLoginPassword(method, uri, data, login, password string) (*http.Request, error) {
	request, err := buildRequest(method, uri, data)
	if err != nil {
		return request, err
	}
	request.SetBasicAuth(login,password)
	return request, nil
}

func buildRequest(method, uri, data string) (*http.Request, error) {
	var requestData *strings.Reader
	if data != "" {
		requestData = strings.NewReader(data)
	} else {
		requestData = strings.NewReader(" ")
		//requestData = nil
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