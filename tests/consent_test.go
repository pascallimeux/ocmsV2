package tests

import (
	"testing"
	"fmt"
	"time"

)
const (
	APPID1      = "APP4TESTS1"
	APPID2      = "APP4TESTS2"
	OWNERID1    = "owner1"
	CONSUMERID1 = "consumer1"
	DATATYPE1   = "type1"
	DATAACCESS1 = "access1"
)

func TestGetVersion(t *testing.T) {
	value, err := ConsHelper.GetVersion(CHAINCODEID)
	if err != nil {
		t.Error("GetVersion return error: ", err)
	}
	fmt.Printf("*** Version of consentV2 CC: %s\n", value)
}

func TestCreateConsent(t *testing.T) {
	txID, err := ConsHelper.CreateConsent(CHAINCODEID, APPID1, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	fmt.Printf("*** TxID for consent creation: %s\n", txID)
}

func TestGetConsents(t *testing.T) {
	consents, err := ConsHelper.GetConsents(CHAINCODEID, APPID1)
	if err != nil {
		fmt.Errorf("GetConsents return error: ", err)
	}
	fmt.Printf("*** GetConsents of consentV2 CC: %s\n", consents)
}

func TestGetConsent(t *testing.T) {
	consentID, err := ConsHelper.CreateConsent(CHAINCODEID, APPID2, OWNERID1, CONSUMERID1, DATATYPE1, DATAACCESS1, getStringDateNow(0), getStringDateNow(7))
	if err != nil {
		t.Error("CreateConsent return error: ", err)
	}
	consent, err := ConsHelper.GetConsent(CHAINCODEID, APPID2, consentID)
	if err != nil {
		fmt.Errorf("GetConsent return error: ", err)
	}
	fmt.Printf("*** GetConsent of consentV2 CC: %s\n", consent)
}









func getStringDateNow(nbdaysafter time.Duration) string{
	t := time.Now().Add(nbdaysafter * 24 * time.Hour)
	return t.Format("2006-01-02")
}