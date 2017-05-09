package helpers
import (
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
	//"errors"
)

type ConsentHelper struct {
	ChainID         string
	Chain 	        fabricClient.Chain
	EventHub        events.EventHub
}

type Consent struct {
	Action		string     `json:"action"`
	AppID 		string     `json:"appid"`
	State       	string     `json:"state"`
	ConsentID      	string     `json:"consentid"`
	OwnerID       	string     `json:"ownerid"`
	ConsumerID      string     `json:"consumerid"`
	DataType      	string     `json:"datatype"`
	DataAccess      string     `json:"dataaccess"`
	Dt_begin      	string     `json:"dtbegin"`
	Dt_end       	string     `json:"dtend"`
}

func (c *Consent) Print() string {
	consentStr := fmt.Sprintf("ConsentID:%s ConsumerID:%s OwnerID:%s Datatype:%s Dataaccess:%s Dt_begin:%s Dt_end:%s", c.ConsentID, c.ConsumerID, c.OwnerID, c.DataType, c.DataAccess, c.Dt_begin, c.Dt_end)
	return consentStr
}


func (ch *ConsentHelper) GetVersion(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "getversion")
	return ch.query(chainCodeID, args)
}

func (ch *ConsentHelper) GetConsent(chainCodeID, appID, consentID string) (Consent, error) {
	var args []string
	args = append(args, "getconsent")
	args = append(args, appID)
	args = append(args, consentID)
	strResp, err := ch.query(chainCodeID, args)
	return extractConsent(consentID, strResp, err)
}

func (ch *ConsentHelper) GetConsents(chainCodeID, appID string) ([]Consent, error) {
	var args []string
	args = append(args, "getconsents")
	args = append(args, appID)
	return extractConsents(ch.query(chainCodeID, args))
}

func (ch *ConsentHelper) GetOwnerConsents(chainCodeID, appID, ownerID string) ([]Consent, error) {
	var args []string
	args = append(args, "getownerconsents")
	args = append(args, appID)
	args = append(args, ownerID)
	return extractConsents(ch.query(chainCodeID, args))
}

func (ch *ConsentHelper) GetConsumerConsents(chainCodeID, appID, consumerID string) ([]Consent, error) {
	var args []string
	args = append(args, "getconsumerconsents")
	args = append(args, appID)
	args = append(args, consumerID)
	return extractConsents(ch.query(chainCodeID, args))
}

func (ch *ConsentHelper) CreateConsent(chainCodeID, appID, ownerID, consumerID, datatype, dataaccess, st_date, end_date string) (string, error) {
	var args []string
	args = append(args, "postconsent")
	args = append(args, appID)
	args = append(args, ownerID)
	args = append(args, consumerID)
	args = append(args, datatype)
	args = append(args, dataaccess)
	args = append(args, st_date)
	args = append(args, end_date)
	txID, err := ch.createTransaction(chainCodeID, args)
	return txID, err
}

func (ch *ConsentHelper) DeleteConsents4Application(chainCodeID, appID string) (string, error) {
	var args []string
	args = append(args, "resetconsents")
	args = append(args, appID)
	txID, err := ch.createTransaction(chainCodeID, args)
	return txID, err
}

func (ch *ConsentHelper) RemoveConsent(chainCodeID, appID, consentID string) (string, error) {
	var args []string
	args = append(args, "removeconsent")
	args = append(args, appID)
	args = append(args, consentID)
	txID, err := ch.createTransaction(chainCodeID, args)
	return txID, err
}

func (ch *ConsentHelper) IsConsentExist(chainCodeID, appID, ownerID, consumerID, dataType, dataAccess string) (bool, error) {
	var args []string
	args = append(args, "isconsent")
	args = append(args, appID)
	args = append(args, ownerID)
	args = append(args, consumerID)
	args = append(args, dataType)
	args = append(args, dataAccess)
	return extractIsConsent(ch.query(chainCodeID, args))
}

func (ch *ConsentHelper) CreateConsentWithRegistration(chainCodeID, appID, ownerID, consumerID, datatype, dataaccess, st_date, end_date string) (string, error) {
	var args []string
	args = append(args, "postconsent")
	args = append(args, appID)
	args = append(args, ownerID)
	args = append(args, consumerID)
	args = append(args, datatype)
	args = append(args, dataaccess)
	args = append(args, st_date)
	args = append(args, end_date)
	txID, err := ch.createTransactionWithRegistration(chainCodeID, args)
	return txID, err
}

func (ch *ConsentHelper) query(chainCodeID string, args []string) (string, error) {
	log.Debug("query(chainCodeID:"+ chainCodeID+" args:"+ strings.Join(args," ") +") : calling method -")
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponses, _, err := sdkUtil.CreateAndSendTransactionProposal(ch.Chain, chainCodeID, ch.ChainID, args, []fabricClient.Peer{ch.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		log.Error("CreateAndSendTransactionProposal return error: %v", err)
		return "", fmt.Errorf("Query CC return error")
	}
	response := string(transactionProposalResponses[0].GetResponsePayload())
	return response, nil
}

func (ch *ConsentHelper) createTransaction(chainCodeID string, args []string) (string, error) {
	log.Debug("createTransaction(chainCodeID:"+ chainCodeID+" args:"+ strings.Join(args," ") +") : calling method -")
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponse, txID, err := sdkUtil.CreateAndSendTransactionProposal(ch.Chain, chainCodeID, ch.ChainID, args, []fabricClient.Peer{ch.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		log.Error("CreateAndSendTransactionProposal return error: %v", err)
		return "", fmt.Errorf("CreateTransactionProposal for CC return error")
	}
	_, err = sdkUtil.CreateAndSendTransaction(ch.Chain, transactionProposalResponse)
	if err != nil {
		log.Error("CreateAndSendTransaction return error: %v", err)
		return "", fmt.Errorf("CreateTransaction for CC return error")
	}
	return txID, nil
}

func (ch *ConsentHelper) createTransactionWithRegistration(chainCodeID string, args []string) (string, error) {
	log.Debug("createTransactionWithRegistration(chainCodeID:"+ chainCodeID+" args:"+ strings.Join(args," ") +") : calling method -")
	eventID := "test([a-zA-Z]+)"
	// Register callback for chaincode event
	done1, rce := sdkUtil.RegisterCCEvent(chainCodeID, eventID, ch.EventHub)
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponse, txID, err := sdkUtil.CreateAndSendTransactionProposal(ch.Chain, chainCodeID, ch.ChainID, args, []fabricClient.Peer{ch.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	// Register for commit event
	done, fail := sdkUtil.RegisterTxEvent(txID, ch.EventHub)

	_, err = sdkUtil.CreateAndSendTransaction(ch.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}

	select {
	case <-done:
	case <-fail:
		return txID, fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return txID, fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	select {
	case <-done1:
	case <-time.After(time.Second * 20):
		return txID, fmt.Errorf("Did NOT receive CC for eventId(%s)\n", eventID)
	}
	ch.EventHub.UnregisterChaincodeEvent(rce)

	return txID, nil
}


func extractConsents(stringresp string, err error) ([]Consent, error) {
	var consents []Consent
	if err != nil {
		return consents, err
	}
	dec := json.NewDecoder(strings.NewReader(stringresp))
	err = dec.Decode(&consents)
	if err != nil {
		log.Error(err)
		err = fmt.Errorf("Extract consents return error")
	}
	return consents, err
}

func extractConsent(consentID, stringresp string, err error) (Consent, error) {
	var consent Consent
	if err != nil {
		return consent, err
	}
	dec := json.NewDecoder(strings.NewReader(stringresp))
	err = dec.Decode(&consent)
	if err != nil {
		log.Error(err)
		err = fmt.Errorf("Extract a consent return error")
	}
	return consent, err
}

func extractIsConsent(stringresp string, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	if stringresp == "True" {
		return true, nil
	} else {
		return false, nil
	}
}