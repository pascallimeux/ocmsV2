package api

import (
	"github.com/pascallimeux/ocmsV2/helpers"
	"testing"
	"errors"
	"net/http"
	"encoding/json"
	"github.com/hyperledger/fabric/protos/common"
)

func TestBlockchainInfo(t *testing.T) {
	var response *common.BlockchainInfo
	request, err := buildRequestWithLoginPassword("GET", httpServerTest.URL+BCINFO, "", ADMINNAME, ADMINPWD)
	if err != nil {
		t.Error(err)
	}
	status, body_bytes, err := executeRequest(request)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(body_bytes, &response)
	if err != nil {
		t.Error(err)
	}
	if status != http.StatusOK {
		t.Error(errors.New("bad status"))
	}
	if response.CurrentBlockHash == nil || response.PreviousBlockHash == nil || response.Height<0{
		t.Error(errors.New("bad response"))
	}
}
