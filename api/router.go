package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("ocms.api")

const (
	VERSIONURI       = "/ocms/v2/api/version"
	CONSENTAPI       = "/ocms/v2/api/consent/"

	BCINFO           = "/ocms/v2/dashboard/chain"
	QUERYTRANSACTION = "/ocms/v2/dashboard/transaction"
	BLOCKBYNB        = "/ocms/v2/dashboard/blocks/nb"
	BLOCKBYHASH      = "/ocms/v2/dashboard/blocks/hash"
	GETCHANNELS      = "/ocms/v2/dashboard/channels"
	INSTALLEDCC      = "/ocms/v2/dashboard/cc/installed"
	QUERYBYCC        = "/ocms/v2/dashboard/cc/query"
	GETPEERS         = "/ocms/v2/dashboard/peers"
	INSTANCIATEDCC   = "/ocms/v2/dashboard/cc/instanciated"

	REGISTER         = "/ocms/v2/admin/user/register"
	ENROLL           = "/ocms/v2/admin/user/enroll"
	REVOKE           = "/ocms/v2/admin/user/revoke"
)

type AppContext struct {
	HttpServer     	*http.Server
	ChainCodeID   	string
	ChainID         string
	Repo            string
	StatStorePath   string
}

func (a *AppContext) CreateOCMSRoutes(router *mux.Router) {
	log.Debug("CreateOCMSRoutes() : calling method -")
	router.HandleFunc(VERSIONURI, a.getVersion).Methods("GET")
	router.HandleFunc(CONSENTAPI, a.processConsent).Methods("POST")
	router.HandleFunc(BCINFO, a.blockchainInfo).Methods("GET")
	router.HandleFunc(GETCHANNELS, a.getChannels).Methods("GET")
	router.HandleFunc(GETPEERS, a.getPeers).Methods("GET")
	router.HandleFunc(INSTALLEDCC, a.getInstalledCC).Methods("GET")
	router.HandleFunc(INSTANCIATEDCC, a.getInstantiatedCC).Methods("GET")
	router.HandleFunc(QUERYTRANSACTION+"/{truuid}", a.transactionDetails).Methods("GET")
	router.HandleFunc(BLOCKBYNB+"/{blocknb}", a.blockByNumber).Methods("GET")
	router.HandleFunc(BLOCKBYHASH+"/{blockhash}", a.blockByHash).Methods("GET")
	router.HandleFunc(QUERYBYCC+"/{ccname}", a.queryByCC).Methods("GET")
	router.HandleFunc(REGISTER, a.registerUser).Methods("POST")
	router.HandleFunc(ENROLL, a.enrollUser).Methods("POST")
	router.HandleFunc(REVOKE, a.revokeUser).Methods("POST")
}
