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
var log = logging.MustGetLogger("api")

const (
	VERSIONURI   = "/ocms/v2/api/version"
	CONSENTAPI   = "/ocms/v2/api/consent/"
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
}
