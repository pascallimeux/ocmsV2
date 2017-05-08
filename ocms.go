package main

import(
	"github.com/pascallimeux/ocmsV2/helpers"
	"github.com/cloudflare/cfssl/log"
	"github.com/pascallimeux/ocmsV2/api"
	"net/http"
	"time"
	"github.com/op/go-logging"
	"github.com/gorilla/mux"
)

var log = logging.MustGetLogger("ocms")


const(
	chainCodePath    = "github.com/consentv2"
	chainCodeVersion = "v0"
	chainCodeID      = "consentv2"
	adminusername    = "admin"
	adminPwd         = "adminpw"
	repo             = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures"
	statstorePath    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/enroll_user"
	configfile       = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/config/config.yaml"
	channelConfig    = "/opt/gopath/src/github.com/pascallimeux/ocmsV2/fixtures/channel/testchannel.tx"
	chainID          = "testchannel"
	providerName     = "SW"
)

func main() {

	networkHelper := helpers.NetworkHelper{
		Repo:                   repo,
		ConfigFile:      	configfile,
		ChannelConfig:   	channelConfig,
		ChainID:         	chainID,
	}

	err := networkHelper.InitNetwork(adminusername, adminPwd, statstorePath, providerName)
	if err != nil {
		log.Fatal(err)
	}

	err = networkHelper.DeployCC(chainCodePath, chainCodeVersion, chainCodeID)
	if err != nil {
		log.Fatal(err)
	}

	consentHelper := helpers.ConsentHelper{
		ChainID:         	chainID,
		Chain:			networkHelper.Chain,
		EventHub:		networkHelper.EventHub,
	}

	// Init application context
	appContext := api.AppContext{ConsentHelper: consentHelper, NetworkHelper: networkHelper}

	// Init routes for application
	var router *mux.Router
	appContext.CreateOCMSRoutes(router)

	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Nanosecond,
		WriteTimeout: configuration.WriteTimeout * time.Nanosecond,
	}
	log.Fatal(s.ListenAndServe().Error())

}
