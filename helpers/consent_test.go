package helpers

import (
	"testing"
	"fmt"
	"time"
	"strconv"
)


func TestGetVersion(t *testing.T) {
	value, err := consHelper.GetVersion(configuration.ChainCodeID)
	if err != nil {
		t.Error("GetVersion return error: ", err)
	}
	if value != CCVERSION{
		t.Error(value, " is not equal with ", CCVERSION)
	}
}

func TestCreateConsent(t *testing.T) {
	_, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID1, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
}

func TestGetConsents(t *testing.T) {
	_, err := consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID2)
	if err != nil {
		fmt.Errorf("DeleteConsents4Application return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID2, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID2, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := consHelper.GetConsents(configuration.ChainCodeID, APPID2)
	if err != nil {
		t.Error("GetConsents return error: ", err)
	}
	if len(consents) != 2 {
		t.Error(" Does not get the right number of consents 2 expected, but ", strconv.Itoa(len(consents)), " received...")
	}

}

func TestGetAConsent(t *testing.T) {
	consentID, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID3, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	_, err = consHelper.GetConsent(configuration.ChainCodeID, APPID3, consentID)
	if err != nil {
		t.Error("GetConsent return error: ", err)
	}
}


func TestRemoveAConsent(t *testing.T) {
	consentID, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID3, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	_, err = consHelper.RemoveConsent(configuration.ChainCodeID, APPID3, consentID)
	if err != nil {
		t.Error("RemoveConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	_, err = consHelper.GetConsent(configuration.ChainCodeID, APPID3, consentID)
	if err == nil {
		t.Error("RemoveConsent return error: ", err)
	}
}

func TestGetOwnerConsents(t *testing.T) {
	_, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID4, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID4, OWNERID1, CONSUMERID2, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID4, OWNERID1, CONSUMERID3, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := consHelper.GetOwnerConsents(configuration.ChainCodeID, APPID4, OWNERID1)
	if err != nil {
		t.Error("GetOwnerConsents return error: ", err)
	}
	if len(consents) != 3 {
		t.Error(" Does not get the right number of consents... ")
	}
}

func TestGetConsumerConsents(t *testing.T) {
	_, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID5, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID5, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID5, OWNERID3, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := consHelper.GetConsumerConsents(configuration.ChainCodeID, APPID5, CONSUMERID1)
	if err != nil {
		t.Error("GetConsumerConsents return error: ", err)
	}
	if len(consents) != 3 {
		t.Error(" Does not get the right number of consents")
	}
}

func TestDeleteConsents4Application(t *testing.T) {
	_, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID6, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(configuration.ChainCodeID, APPID6, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	_, err = consHelper.DeleteConsents4Application(configuration.ChainCodeID, APPID6)
	if err != nil {
		t.Error("DeleteConsents4Application return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := consHelper.GetConsents(configuration.ChainCodeID, APPID6)
	if err != nil {
		t.Error("GetConsents return error: ", err)
	}
	if len(consents) != 0 {
		t.Error(" Error not empty list of consents... ")
	}
}

func TestIsConsentExist(t *testing.T) {
	_, err := consHelper.CreateConsent(configuration.ChainCodeID, APPID1, OWNERID3, CONSUMERID3, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(TransactionTimeout)
	exist, err := consHelper.IsConsentExist(configuration.ChainCodeID, APPID1, OWNERID3, CONSUMERID3, DATATYPE1, DATAACCESS1)
	if err != nil {
		t.Error("IsConsentExist return error: ", err)
	}
	if ! exist {
		t.Error("bad response for isConsentExist...")
	}
}

func getStringDateNow(nbdaysafter time.Duration) string{
	t := time.Now().Add(nbdaysafter * 24 * time.Hour)
	return t.Format("2006-01-02")
}