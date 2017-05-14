package helpers
import(
	"fmt"
	sdkUtil "github.com/hyperledger/fabric-sdk-go/fabric-client/helpers"
	sdkConfig "github.com/hyperledger/fabric-sdk-go/config"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"errors"
)

var log = logging.MustGetLogger("ocms.helpers")

type Helper interface {
	Init(UserCredentials) error
}

type NetworkHelper struct {
	ChainID         string
	StatStorePath   string
	Repo            string
	EventHub        events.EventHub
	Client          fabricClient.Client
	Chain 	        fabricClient.Chain
	Initialized	bool
}

func (nh *NetworkHelper) Init(userCredentials UserCredentials) error{
	chain, err := getChain(userCredentials, nh.StatStorePath, nh.ChainID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	client, err := getClient(userCredentials, nh.StatStorePath)
	if err != nil {
		return err
	}
	eventHub, err := getEventHub()
	if err != nil {
		return err
	}
	if err := eventHub.Connect(); err != nil {
		return err
	}
	nh.Chain    = chain
	nh.Client   = client
	nh.EventHub = eventHub
	nh.Initialized = true
	return nil
}

func (nh *NetworkHelper) StartNetwork(userCredentials UserCredentials, providerName, netConfigFile, channelConfig string)  error{
	log.Debug("InitNetwork(username:"+ userCredentials.UserName+" providerName:"+ providerName+") : calling method -")
	initError := fmt.Errorf("InitNetwork return error")
	// Init SDK config
	err := sdkConfig.InitConfig(netConfigFile)
	if err != nil {
		log.Error("Failed init sdk config", err)
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
	err = nh.Init(userCredentials)
	if err != nil {
		log.Error("Failed init networkHandler [",err,"]")
		return initError
	}
	// Create and join channel
	if err := sdkUtil.CreateAndJoinChannel(nh.Client, nh.Chain, channelConfig); err != nil {
		log.Error("CreateAndJoinChannel return error: ", err)
		return initError
	}
	log.Debug("Hyperledger network initialized...")
	return nil
}


func (nh *NetworkHelper) DeployCC(chainCodePath, chainCodeVersion, chainCodeID string) error {
	log.Debug("DeployCC(chainCodePath:"+ chainCodePath+" chainCodeVersion:" + chainCodeVersion +" chainCodeID:"+ chainCodeID+") : calling method -")
	if err := nh.InstallCC(chainCodePath, chainCodeVersion, chainCodeID, nil); err != nil {
		return err
	}
	var args []string
	return nh.InstantiateCC(chainCodePath, chainCodeVersion, chainCodeID, args)
}

func (nh *NetworkHelper) InstallCC(chainCodePath, chainCodeVersion, chainCodeID string, chaincodePackage []byte) error {
	if err := sdkUtil.SendInstallCC(nh.Client, nh.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, nh.Chain.GetPeers(), nh.Repo); err != nil {
		log.Error("SendInstallProposal return error: ", err)
		return fmt.Errorf("Install chaincode return error")
	}
	log.Debug("Chaincode "+chainCodeID+" installed...")
	return nil
}

func (nh *NetworkHelper) InstantiateCC(chainCodePath, chainCodeVersion, chainCodeID string, args []string) error {
	if err := sdkUtil.SendInstantiateCC(nh.Chain, chainCodeID, nh.ChainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{nh.Chain.GetPrimaryPeer()}, nh.EventHub); err != nil {
		log.Error("SendInstantiateProposal return error: ", err)
		return fmt.Errorf("Instantiate chaincode return error")
	}
	log.Debug("Chaincode "+chainCodeID+" Instantiate...")
	return nil
}


func (nh *NetworkHelper) QueryInfos()(*common.BlockchainInfo, error){
	log.Debug("QueryInfos() : calling method -")
	return nh.Chain.QueryInfo()
}

func (nh *NetworkHelper) QueryTransaction(transactionID string)(*pb.ProcessedTransaction, error){
	log.Debug("QueryTransaction("+transactionID+") : calling method -")
	processTransaction, err := nh.Chain.QueryTransaction(transactionID)
	if err == nil {
		log.Debug("PAYLOAD:", processTransaction.TransactionEnvelope.String())
	}
	return processTransaction, err
}

func (nh *NetworkHelper) QueryBlockByNumber(stnb string)(*common.Block, error){
	log.Debug("QueryBlockByNumber("+stnb+") : calling method -")
	nb, err :=strconv.Atoi(stnb)
	if err != nil {
		nb = -1
	}
	return nh.Chain.QueryBlock(nb)
}

func (nh *NetworkHelper) QueryBlockByHash(hash string)(*common.Block, error){
	log.Debug("QueryBlockByHash("+hash+") : calling method -")
	return nh.Chain.QueryBlockByHash([]byte(hash))
}

func (nh *NetworkHelper) QueryChannels()(*pb.ChannelQueryResponse, error){
	log.Debug("QueryChannels() : calling method -")
	target := nh.Chain.GetPrimaryPeer()
	return nh.Client.QueryChannels(target)
}

func (nh *NetworkHelper) GetInstalledChainCode()(*pb.ChaincodeQueryResponse, error){
	target := nh.Chain.GetPrimaryPeer()
	log.Debug("QueryInstalledChaincodes("+target.GetURL()+") : calling method -")
	return  nh.Client.QueryInstalledChaincodes(target)
}

func (nh *NetworkHelper) GetInstanciateChainCode()(*pb.ChaincodeQueryResponse, error){
	log.Debug("GetInstanciateChainCode() : calling method -")
	return nh.Chain.QueryInstantiatedChaincodes()
}

func (nh *NetworkHelper) QueryByChainCode(chaincodeName string)([][]byte, error){
	log.Debug("QueryByChaincode("+chaincodeName+") : calling method -")
	targets := nh.Chain.GetPeers()
	return nh.Chain.QueryByChaincode(chaincodeName, []string{"getinstalledchaincodes"}, targets)
}

func (nh *NetworkHelper) GetPeers()([]fabricClient.Peer){
	log.Debug("GetPeers() : calling method -")
	return nh.Chain.GetPeers()
}

func getChain(userCredentials UserCredentials, statStorePath, chainID string) (fabricClient.Chain, error) {
	log.Debug("GetChain(username:"+ userCredentials.UserName+") : calling method -")
	var chain fabricClient.Chain

	client, err := getClient(userCredentials, statStorePath)
	if err != nil {
		return chain, errors.New("getClient return error: %v" + err.Error())
	}
	chain, err = sdkUtil.GetChain(client, chainID)
	if err != nil {
		log.Error("Create chain ", chainID," failed: ", err)
		return chain, err
	}
	return chain, nil
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
			log.Debug("EventHub connect to peer (", p.EventHost,":", p.EventPort,")")
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
