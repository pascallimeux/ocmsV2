package test

import (
	"testing"
	"github.com/pascallimeux/ocmsV2/ocms"
	"github.com/cloudflare/cfssl/log"
)

func TestChainCodeInvoke(t *testing.T) {

	app := ocms.OcmsApp{
		ConfigFile:      	REPO+"/config/config.yaml",
		ChannelConfig:   	REPO+"/channel/testchannel.tx",
		ChainID:         	"testchannel",
		ConnectEventHub: true,
	}

	err := app.initConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = app.setup()
	if err != nil {
		log.Fatal(err)
	}

}
