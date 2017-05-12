package api
import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocmsV2/helpers"
	"strings"
)

func SendError(w http.ResponseWriter, err error) {
	log.Debug("sendError() : calling method -")
	libelle := err.Error()
	libelle = strings.Replace(libelle, "\"", "'", -1)
	log.Error("sendError: ", libelle)
	message := "{\"content\":\"" + libelle + "\"} "
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

//HTTP Get - /ocms/v2/dashboard/chain
func (a *AppContext) blockchainInfo(w http.ResponseWriter, r *http.Request) {
	log.Debug("blockchainInfo() : calling method -")
	var err error
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err = InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	blockchainInfo, err := netHelper.QueryInfos()
	if err != nil {
		SendError(w, err)
	}
	content, err := json.Marshal(blockchainInfo)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/channels
func (a *AppContext) getChannels(w http.ResponseWriter, r *http.Request) {
	log.Debug("getChannels() : calling method -")
	var err error
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err = InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	channels, err := netHelper.QueryChannels()
	if err != nil {
		SendError(w, err)
	}
	content, err := json.Marshal(channels)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/peers
func (a *AppContext) getPeers(w http.ResponseWriter, r *http.Request) {
	log.Debug("getChannels() : calling method -")
	var err error
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err = InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	peers :=netHelper.GetPeers()
	content, err := json.Marshal(peers)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/cc/installed
func (a *AppContext) getInstalledCC(w http.ResponseWriter, r *http.Request) {
	log.Debug("getInstalledCC() : calling method -")
	var err error
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err = InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	cc, err := netHelper.GetInstalledChainCode()
	if err != nil {
		SendError(w, err)
	}
	content, err := json.Marshal(cc)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/cc/instanciated
func (a *AppContext) getInstantiatedCC(w http.ResponseWriter, r *http.Request) {
	log.Debug("getInstantiatedCC() : calling method -")
	var err error
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err = InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	cc, err := netHelper.GetInstanciateChainCode()
	if err != nil {
		SendError(w, err)
	}
	content, err := json.Marshal(cc)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/transaction/{truuid}
func (a *AppContext) transactionDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tr_uuid := vars["truuid"]
	message := fmt.Sprintf("transactionDetails(tr_uuid=%s) : calling method -", tr_uuid)
	log.Debug(message)
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err := InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	transaction, err := netHelper.QueryTransaction(tr_uuid)
	if err != nil {
		SendError(w, err)
		return
	}
	content, err := json.Marshal(transaction)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

//HTTP Get - /ocms/v2/dashboard/blocks/nb/{blocknb}
func (a *AppContext) blockByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blocNb := vars["blocknb"]
	message := fmt.Sprintf("blockByNumber(blocknb=%s) : calling method -", blocNb)
	log.Debug(message)
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err := InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	block, err := netHelper.QueryBlockByNumber(blocNb)
	if err != nil {
		SendError(w, err)
		return
	}
	content, err := json.Marshal(block)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}


//HTTP Get - /ocms/v2/dashboard/blocks/hash/{blockhash}
func (a *AppContext) blockByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockHash := vars["blockhash"]
	message := fmt.Sprintf("transactionDetails(blockHash=%s) : calling method -", blockHash)
	log.Debug(message)
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err := InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	block, err := netHelper.QueryBlockByHash(blockHash)
	if err != nil {
		SendError(w, err)
		return
	}
	content, err := json.Marshal(block)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}


//HTTP Get - /ocms/v2/dashboard/blocks/hash/{ccname}
func (a *AppContext) queryByCC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chaincodeName := vars["ccname"]
	message := fmt.Sprintf("queryByCC(blockHash=%s) : calling method -", chaincodeName)
	log.Debug(message)
	netHelper := &helpers.NetworkHelper{Repo: a.Repo, StatStorePath: a.StatStorePath, ChainID: a.ChainID}
	err := InitHelper(r, netHelper)
	if err != nil {
		SendError(w, err)
		return
	}
	response, err := netHelper.QueryByChainCode(chaincodeName)
	if err != nil {
		SendError(w, err)
		return
	}
	content, err := json.Marshal(response)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}