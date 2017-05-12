package main

import(
	"github.com/pascallimeux/ocmsV2/helpers"
	"github.com/pascallimeux/ocmsV2/api"
	"github.com/pascallimeux/ocmsV2/settings"
	"net/http"
	"time"
	"github.com/op/go-logging"
	"github.com/gorilla/mux"
)

var log = logging.MustGetLogger("ocms")

func main() {

	// Init settings
	configuration, err := settings.GetSettings(".", "ocms")
	if err != nil {
		panic(err.Error())
	}
	adminCredentials := helpers.UserCredentials {
		UserName:configuration.Adminusername,
		EnrollmentSecret:configuration.AdminPwd}

	// Init Hyperledger network
	networkHelper := helpers.NetworkHelper{
		Repo:                   configuration.Repo,
		StatStorePath:          configuration.StatstorePath,
		ChainID:         	configuration.ChainID}
	err = networkHelper.StartNetwork(adminCredentials, configuration.ProviderName, configuration.SDKConfigfile, configuration.ChannelConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy the consent smartcontract if is not deployed
	networkHelper.DeployCC(configuration.ChainCodePath, configuration.ChainCodeVersion, configuration.ChainCodeID)

	// Init application context
	appContext := api.AppContext{
		ChainCodeID: 		configuration.ChainCodeID,
		Repo:                   configuration.Repo,
		StatStorePath:          configuration.StatstorePath,
		ChainID:         	configuration.ChainID,
	}

	// Init routes for application
	router := mux.NewRouter().StrictSlash(false)
	appContext.CreateOCMSRoutes(router)

	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Nanosecond,
		WriteTimeout: configuration.WriteTimeout * time.Nanosecond,
	}
	log.Fatal(s.ListenAndServe().Error())

	defer configuration.Close()

}
