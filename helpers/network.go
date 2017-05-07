package helpers
import(
	"fmt"
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	sdkConfig "github.com/hyperledger/fabric-sdk-go/config"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("helpers")

type NetworkHelper struct {
	ChainID         string
	Repo            string
	ConfigFile     	string
	ChannelConfig	string
	Client          fabricClient.Client
	Chain 	        fabricClient.Chain
	EventHub        events.EventHub
	Initialized     bool
}

func (nh *NetworkHelper) InitNetwork(username, password, stateStorePath, providerName string)  error{
	log.Debug("InitNetwork(username:"+ username+" stateStorePath:"+ stateStorePath +" providerName:"+ providerName+") : calling method -")
	initError := fmt.Errorf("InitNetwork return error")
	// Init config
	err := sdkConfig.InitConfig(nh.ConfigFile)
	if err != nil {
		log.Error("Failed init sdk-go config", err)
		return initError
	}
	err = bccspFactory.InitFactories(&bccspFactory.FactoryOpts{
		ProviderName: providerName,
		SwOpts: &bccspFactory.SwOpts{
			HashFamily: sdkConfig.GetSecurityAlgorithm(),
			SecLevel:   sdkConfig.GetSecurityLevel(),
			FileKeystore: &bccspFactory.FileKeystoreOpts{
				KeyStorePath: sdkConfig.GetKeyStorePath(),
			},
			Ephemeral: false,
		},
	})
	if err != nil {
		log.Error("Failed getting ephemeral software-based BCCSP [",err,"]")
		return initError
	}
	// Get client
	client, err := sdkUtil.GetClient(username, password, stateStorePath)
	if err != nil {
		log.Error("Create client failed: ", err)
		return initError
	}
	nh.Client = client

	// Get chain
	chain, err := sdkUtil.GetChain(client, nh.ChainID)
	if err != nil {
		log.Error("Create chain ", nh.ChainID," failed: ", err)
		return initError
	}
	nh.Chain = chain

	// Create and join channel
	if err := sdkUtil.CreateAndJoinChannel(nh.Client, nh.Chain, nh.ChannelConfig); err != nil {
		log.Error("CreateAndJoinChannel return error: ", err)
		return initError
	}

	// Get envenHub
	eventHub, err := getEventHub()
	if err != nil {
		log.Error("Fail get eventHub: ", err)
		return initError
	}

	if err := eventHub.Connect(); err != nil {
		log.Error("Failed eventHub.Connect() ", err)
		return initError
	}

	nh.EventHub = eventHub
	nh.Initialized = true
	log.Debug("Hyperledger network initialized...")
	return nil
}


func (nh *NetworkHelper) DeployCC(chainCodePath, chainCodeVersion, chainCodeID string) error {
	log.Debug("DeployCC(chainCodePath:"+ chainCodePath+" chainCodeVersion:" + chainCodeVersion +" chainCodeID:"+ chainCodeID+") : calling method -")
	if err := nh.installCC(chainCodePath, chainCodeVersion, chainCodeID, nil); err != nil {
		return err
	}
	var args []string
	return nh.instantiateCC(chainCodePath, chainCodeVersion, chainCodeID, args)
}

func (nh *NetworkHelper) installCC(chainCodePath, chainCodeVersion, chainCodeID string, chaincodePackage []byte) error {
	if err := sdkUtil.SendInstallCC(nh.Client, nh.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, nh.Chain.GetPeers(), nh.Repo); err != nil {
		log.Error("SendInstallProposal return error: ", err)
		return fmt.Errorf("Install chaincode return error")
	}
	log.Debug("Chaincode "+chainCodeID+" installed...")
	return nil
}

func (nh *NetworkHelper) instantiateCC(chainCodePath, chainCodeVersion, chainCodeID string, args []string) error {
	if err := sdkUtil.SendInstantiateCC(nh.Chain, chainCodeID, nh.ChainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{nh.Chain.GetPrimaryPeer()}, nh.EventHub); err != nil {
		log.Error("SendInstantiateProposal return error: ", err)
		return fmt.Errorf("Instantiate chaincode return error")
	}
	log.Debug("Chaincode "+chainCodeID+" Instantiate...")
	return nil
}

func getEventHub() (events.EventHub, error) {
	log.Debug("getEventHub() : calling method -")
	eventHub := events.NewEventHub()
	foundEventHub := false
	peerConfig, err := sdkConfig.GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading peer config: %v", err)
	}
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("******* EventHub connect to peer (%s:%d) *******\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, fmt.Errorf("No EventHub configuration found")
	}
	return eventHub, nil
}
