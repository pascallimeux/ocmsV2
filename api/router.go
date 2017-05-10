/*
Copyright Pascal Limeux. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocmsV2/helpers"
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
)

type AppContext struct {
	HttpServer     	*http.Server
	ConsentHelper	helpers.ConsentHelper
	NetworkHelper 	helpers.NetworkHelper
	ChainCodeID   	string
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
	router.HandleFunc(QUERYTRANSACTION+"/{truuid}", a.queryByCC).Methods("GET")
	router.HandleFunc(BLOCKBYNB+"/{blocknb}", a.blockByNumber).Methods("GET")
	router.HandleFunc(BLOCKBYHASH+"/{blockhash}", a.blockByHash).Methods("GET")
	router.HandleFunc(QUERYBYCC+"/{ccname}", a.queryByCC).Methods("GET")
}
