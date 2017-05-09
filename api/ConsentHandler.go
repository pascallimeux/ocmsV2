package api

import (
	"net/http"
	"encoding/json"
	"time"
	"fmt"
	"errors"
	"github.com/pascallimeux/ocmsV2/helpers"
)


type Version struct {
	Version string
}

type IsConsent struct {
	Consent string
}


//HTTP Get - /ocms/v2/api/version
func (a *AppContext) getVersion(w http.ResponseWriter, r *http.Request) {
	log.Debug("getVersion() : calling method -")
	version := Version{}
	var err error
	version.Version, err = a.ConsentHelper.GetVersion(a.ChainCodeID)
	if err != nil {
		SendError(w, err)
	}
	content, _ := json.Marshal(version)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}


//HTTP Post - /ocms/v2/api/consent
func (a *AppContext) processConsent(w http.ResponseWriter, r *http.Request) {
	log.Debug("processConsent() : calling method -")

	var bytes []byte
	var consent helpers.Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		SendError(w, err)
		return
	}
	switch action := consent.Action; action {
	case "create":
		bytes, err = a.createConsent(a.ChainCodeID, consent)
	case "list":
		bytes, err = a.listConsents(a.ChainCodeID, consent.AppID)
	case "get":
		bytes, err = a.getConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "remove":
		bytes, err = a.unactivateConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "list4owner":
		bytes, err = a.getConsents4Owner(a.ChainCodeID, consent.AppID, consent.OwnerID)
	case "list4consumer":
		bytes, err = a.getConsents4Consumer(a.ChainCodeID, consent.AppID, consent.ConsumerID)
	case "isconsent":
		bytes, err = a.isConsent(a.ChainCodeID, consent)
	default:
		log.Error("bad action request")
		SendError(w, err)
		return
	}
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}


func (a *AppContext) createConsent(chainCodeID string, consent helpers.Consent) ([]byte, error) {
	err := check_args(&consent)
	var message string
	if err != nil {
		message = fmt.Sprintf("createConsent(%s) : calling method -", err.Error())
	} else {
		message = fmt.Sprintf("createConsent(%s) : calling method -", consent.Print())
	}
	log.Info(message)
	if err != nil {
		return nil, err
	}
	consentID, err := a.ConsentHelper.CreateConsent(chainCodeID, consent.AppID, consent.OwnerID, consent.ConsumerID, consent.DataType, consent.DataAccess, consent.Dt_begin, consent.Dt_end)
	if err != nil {
		return nil, err
	}
	consent.ConsentID = consentID
	return consent2Bytes(consent)
}

func (a *AppContext) listConsents(chainCodeID, applicationID string) ([]byte, error) {
	message := fmt.Sprintf("listConsents(applicationID=%s) : calling method -", applicationID)
	log.Info(message)
	consents, err := a.ConsentHelper.GetConsents(chainCodeID, applicationID)
	if err != nil {
		return nil, err
	}
	return consents2Bytes(consents)
}

func (a *AppContext) getConsent(chainCodeID, applicationID, consentID string) ([]byte, error) {
	message := fmt.Sprintf("getConsent(applicationID=%s, consentID=%s) : calling method -", applicationID, consentID)
	log.Info(message)
	consent, err := a.ConsentHelper.GetConsent(chainCodeID, applicationID, consentID)
	if err != nil {
		return nil, err
	}
	return consent2Bytes(consent)
}

func (a *AppContext) unactivateConsent(chainCodeID, applicationID, consentID string) ([]byte, error) {
	message := fmt.Sprintf("unactivateConsent(applicationID=%s, consentID=%s) : calling method -", applicationID, consentID)
	log.Info(message)
	_, err := a.ConsentHelper.RemoveConsent(chainCodeID, applicationID, consentID)
	if err != nil {
		return nil, err
	}
	consent, err := a.ConsentHelper.GetConsent(chainCodeID, applicationID, consentID)
	if err != nil {
		return nil, err
	}
	return consent2Bytes(consent)
}

func (a *AppContext) getConsents4Consumer(chainCodeID, applicationID, consumerID string) ([]byte, error) {
	message := fmt.Sprintf("getConsents4Consumer(applicationID=%s, consumerID=%s) : calling method -", applicationID, consumerID)
	log.Info(message)
	consents, err := a.ConsentHelper.GetConsumerConsents(chainCodeID, applicationID, consumerID)
	if err != nil {
		return nil, err
	}
	return consents2Bytes(consents)
}

func (a *AppContext) getConsents4Owner(chainCodeID, applicationID, ownerID string) ([]byte, error) {
	message := fmt.Sprintf("getConsents4Owner(applicationID=%s, ownerID=%s) : calling method -", applicationID, ownerID)
	log.Info(message)
	consents, err := a.ConsentHelper.GetOwnerConsents(chainCodeID, applicationID, ownerID)
	if err != nil {
		return nil, err
	}
	return consents2Bytes(consents)
}

func (a *AppContext) isConsent(chainCodeID string, consent helpers.Consent) ([]byte, error) {
	message := fmt.Sprintf("isConsent(consent=%s) : calling method -", consent.Print())
	log.Info(message)
	isconsent, err := a.ConsentHelper.IsConsentExist(chainCodeID, consent.AppID, consent.OwnerID, consent.ConsumerID, consent.DataType, consent.DataAccess)
	if err != nil {
		return nil, err
	}
	response := IsConsent{}
	if isconsent {
		response.Consent = "True"
	} else {
		response.Consent = "False"
	}
	content, _ := json.Marshal(response)
	return content, nil
}

func consents2Bytes(consents []helpers.Consent) ([]byte, error) {
	log.Debug("consents2Bytes() : calling method -")
	j, err := json.Marshal(consents)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func consent2Bytes(consent helpers.Consent) ([]byte, error) {
	log.Debug("consent2Bytes() : calling method -")
	j, err := json.Marshal(consent)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func check_args(consent *helpers.Consent) error {
	log.Debug("check_args() : calling method -")
	if consent.AppID == "" {
		return errors.New("appID is mandatory!")
	}
	if consent.OwnerID == "" {
		return errors.New("ownerID is mandatory!")
	}
	if consent.ConsumerID == "" {
		return errors.New("consumerID is mandatory!")
	}
	if consent.DataAccess == "" {
		consent.DataAccess = "A"
	}
	if consent.DataType == "" {
		consent.DataType = "All"
	}
	if consent.Dt_begin == "" {
		consent.Dt_begin = time.Now().Format("2006-01-02")
	}
	if consent.Dt_end == "" {
		consent.Dt_end = "2099-01-01"
	}
	return nil
}
