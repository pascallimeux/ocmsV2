package helpers
import(
	"fmt"
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	sdkConfig "github.com/hyperledger/fabric-sdk-go/config"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
)

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
	// Init config
	err := sdkConfig.InitConfig(nh.ConfigFile)
	if err != nil {
		return fmt.Errorf("Failed init sdk-go config [%v]", err)
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
		return fmt.Errorf("Failed getting ephemeral software-based BCCSP [%v]", err)
	}
	// Get client
	client, err := sdkUtil.GetClient(username, password, stateStorePath)
	if err != nil {
		return fmt.Errorf("Create client failed: %v", err)
	}
	nh.Client = client

	// Get chain
	chain, err := sdkUtil.GetChain(client, nh.ChainID)
	if err != nil {
		return fmt.Errorf("Create chain (%s) failed: %v", nh.ChainID, err)
	}
	nh.Chain = chain

	// Create and join channel
	if err := sdkUtil.CreateAndJoinChannel(nh.Client, nh.Chain, nh.ChannelConfig); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	// Get envenHub
	eventHub, err := getEventHub()
	if err != nil {
		return err
	}

	if err := eventHub.Connect(); err != nil {
		return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
	}

	nh.EventHub = eventHub
	nh.Initialized = true
	fmt.Println("Hyperledger network initialized...")
	return nil
}


func (nh *NetworkHelper) DeployCC(chainCodePath, chainCodeVersion, chainCodeID string) error {
	if err := nh.installCC(chainCodePath, chainCodeVersion, chainCodeID, nil); err != nil {
		return err
	}
	var args []string
	return nh.instantiateCC(chainCodePath, chainCodeVersion, chainCodeID, args)
}

func (nh *NetworkHelper) installCC(chainCodePath, chainCodeVersion, chainCodeID string, chaincodePackage []byte) error {
	if err := sdkUtil.SendInstallCC(nh.Client, nh.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, nh.Chain.GetPeers(), nh.Repo); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}
	fmt.Println("Chaincode "+chainCodeID+" installed...")
	return nil
}

func (nh *NetworkHelper) instantiateCC(chainCodePath, chainCodeVersion, chainCodeID string, args []string) error {
	if err := sdkUtil.SendInstantiateCC(nh.Chain, chainCodeID, nh.ChainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{nh.Chain.GetPrimaryPeer()}, nh.EventHub); err != nil {
		return err
	}
	fmt.Println("Instantiate OK...")
	return nil
}

func getEventHub() (events.EventHub, error) {
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
