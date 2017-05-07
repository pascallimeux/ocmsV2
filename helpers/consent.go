package helpers
import (
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
)

type ConsentHelper struct {
	ChainID         string
	Chain 	        fabricClient.Chain
	EventHub        events.EventHub
}

func (ch *ConsentHelper) GetVersion(chainCodeID string) (string, error) {
	var args []string
	args = append(args, "getversion")
	return ch.Query(chainCodeID, args)
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
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponse, txID, err := sdkUtil.CreateAndSendTransactionProposal(ch.Chain, chainCodeID, ch.ChainID, args, []fabricClient.Peer{ch.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	_, err = sdkUtil.CreateAndSendTransaction(ch.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	return txID, nil
}

func (ch *ConsentHelper) CreateConsentWithAck(chainCodeID, appID, ownerID, consumerID, datatype, dataaccess, st_date, end_date string) (string, error) {
	eventID := "test([a-zA-Z]+)"

	// Register callback for chaincode event
	done1, rce := sdkUtil.RegisterCCEvent(chainCodeID, eventID, ch.EventHub)

	var args []string
	args = append(args, "postconsent")
	args = append(args, appID)
	args = append(args, ownerID)
	args = append(args, consumerID)
	args = append(args, datatype)
	args = append(args, dataaccess)
	args = append(args, st_date)
	args = append(args, end_date)
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

func (ch *ConsentHelper) GetConsents(chainCodeID, appID string) (string, error) {
	var args []string
	args = append(args, "getconsents")
	args = append(args, appID)
	return ch.Query(chainCodeID, args)
}

func (ch *ConsentHelper) GetConsent(chainCodeID, appID, consentID string) (string, error) {
	var args []string
	args = append(args, "getconsent")
	args = append(args, appID)
	args = append(args, consentID)
	return ch.Query(chainCodeID, args)
}

func (ch *ConsentHelper) Query(chainCodeID string, args []string) (string, error) {
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("TODO change...")
	transactionProposalResponses, _, err := sdkUtil.CreateAndSendTransactionProposal(ch.Chain, chainCodeID, ch.ChainID, args, []fabricClient.Peer{ch.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	return string(transactionProposalResponses[0].GetResponsePayload()), nil
}