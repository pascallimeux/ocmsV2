package tests

import (
	"testing"
	"fmt"
	"time"
	"strconv"
)


func TestGetVersion(t *testing.T) {
	value, err := consHelper.GetVersion(CHAINCODEID)
	if err != nil {
		t.Error("GetVersion return error: ", err)
	}
	if value != CCVERSION{
		t.Error(value, " is not equal with ", CCVERSION)
	}
}

func TestCreateConsent(t *testing.T) {
	_, err := consHelper.CreateConsent(CHAINCODEID, APPID1, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
}

func TestGetConsents(t *testing.T) {
	_, err := consHelper.DeleteConsents4Application(CHAINCODEID, APPID2)
	if err != nil {
		fmt.Errorf("DeleteConsents4Application return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID2, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID2, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	consents, err := consHelper.GetConsents(CHAINCODEID, APPID2)
	if err != nil {
		t.Error("GetConsents return error: ", err)
	}
	if len(consents) != 2 {
		t.Error(" Does not get the right number of consents 2 expected, but ", strconv.Itoa(len(consents)), " received...")
	}

}

func TestGetAConsent(t *testing.T) {
	consentID, err := consHelper.CreateConsent(CHAINCODEID, APPID3, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	_, err = consHelper.GetConsent(CHAINCODEID, APPID3, consentID)
	if err != nil {
		t.Error("GetConsent return error: ", err)
	}
}


func TestRemoveAConsent(t *testing.T) {
	consentID, err := consHelper.CreateConsent(CHAINCODEID, APPID3, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	_, err = consHelper.RemoveConsent(CHAINCODEID, APPID3, consentID)
	if err != nil {
		t.Error("RemoveConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	_, err = consHelper.GetConsent(CHAINCODEID, APPID3, consentID)
	if err == nil {
		t.Error("RemoveConsent return error: ", err)
	}
}

func TestGetOwnerConsents(t *testing.T) {
	_, err := consHelper.CreateConsent(CHAINCODEID, APPID4, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID4, OWNERID1, CONSUMERID2, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID4, OWNERID1, CONSUMERID3, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	consents, err := consHelper.GetOwnerConsents(CHAINCODEID, APPID4, OWNERID1)
	if err != nil {
		t.Error("GetOwnerConsents return error: ", err)
	}
	if len(consents) != 3 {
		t.Error(" Does not get the right number of consents... ")
	}
}

func TestGetConsumerConsents(t *testing.T) {
	_, err := consHelper.CreateConsent(CHAINCODEID, APPID5, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID5, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID5, OWNERID3, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	consents, err := consHelper.GetConsumerConsents(CHAINCODEID, APPID5, CONSUMERID1)
	if err != nil {
		t.Error("GetConsumerConsents return error: ", err)
	}
	if len(consents) != 3 {
		t.Error(" Does not get the right number of consents")
	}
}

func TestDeleteConsents4Application(t *testing.T) {
	_, err := consHelper.CreateConsent(CHAINCODEID, APPID6, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	_, err = consHelper.CreateConsent(CHAINCODEID, APPID6, OWNERID2, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	_, err = consHelper.DeleteConsents4Application(CHAINCODEID, APPID6)
	if err != nil {
		t.Error("DeleteConsents4Application return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	consents, err := consHelper.GetConsents(CHAINCODEID, APPID6)
	if err != nil {
		t.Error("GetConsents return error: ", err)
	}
	if len(consents) != 0 {
		t.Error(" Error not empty list of consents... ")
	}
}

func TestIsConsentExist(t *testing.T) {
	_, err := consHelper.CreateConsent(CHAINCODEID, APPID1, OWNERID3, CONSUMERID3, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	time.Sleep(time.Millisecond * 1500)
	exist, err := consHelper.IsConsentExist(CHAINCODEID, APPID1, OWNERID3, CONSUMERID3, DATATYPE1, DATAACCESS1)
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