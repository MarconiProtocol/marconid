package magent_base

import (
  "../../../core/blockchain"
  "../../../core/blockchain/vars"
  "../../../core/config"
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/core/manager"
  "../../../core/net/dht"
  "../../../core/net/vars"
  "../../../core/peer"
  "../../../core/runtime"
  "../middleware/interface"
  "./peer_updates"
  "crypto/rsa"
  "fmt"
  "github.com/MarconiProtocol/log"
  "github.com/armon/go-socks5"
  "os"
  "strings"
  "sync"
  "time"
)

/*
  AgentClient encapsulates a Marconi service agent
*/
type AgentClient struct {
  teardownSignal                 chan os.Signal
  peerResponseHandlerStatus      map[string]bool
  peerResponseHandlerStatusMutex sync.Mutex
  baseBeaconKey                  *rsa.PrivateKey
  baseL2KeyFilePath              string
  edgeBeaconKey                  *rsa.PrivateKey
  socks5Server                   *socks5.Server
}

/*
  Return a new instance of AgentClient
*/
func NewAgentClient(conf *AgentConfig) *AgentClient {
  agentClient := AgentClient{}
  agentClient.initialize(conf)
  Agent = &agentClient
  return &agentClient
}

/*
  Start the agent daemon, connections will be created when peers are registered to this node
*/
func (agent *AgentClient) Start() {
  keyManager := mcrypto_key.KeyManagerInstance()

  // Wait until there is agent configured network contract address
  agent.idleForNetworkContractAddress()
  // Wait for agent successful registration to the middleware
  agent.idleForMiddlewareRegistration(keyManager.GetBasePublicKeyHash())

  // Start the base route and peer route beacons
  mnet_dht.GetBeaconManager().CreateBaseRouteBeacon()
  mnet_dht.GetBeaconManager().StartBaseRouteAnnouncement(&agent.baseBeaconKey.PublicKey)
  mnet_dht.GetBeaconManager().CreatePeerRouteBeacon(agent.requestPeerResponseHandler)
  mnet_dht.GetBeaconManager().StartPeerRouteAnnouncement()

  // Start the edge route beacon
  if agent.edgeBeaconKey != nil {
    mnet_dht.GetBeaconManager().CreateEdgeRouteBeacon(nil)
    mnet_dht.GetBeaconManager().StartEdgeRouteAnnouncement(&agent.edgeBeaconKey.PublicKey)
  }

  // peerUpdatesChannel is populated when agent relevant peer update is received
  peerUpdatesChannel := mblockchain.GetBlockchainManager().GetPeerUpdates()
  for {
    peerUpdate := <-peerUpdatesChannel
    mlog.GetLogger().Info(fmt.Sprintf("Received PeerUpdate, action: %s, peer: %s", peerUpdate.Action, peerUpdate.PeerPubKeyHash))

    switch peerUpdate.Action {
    case mblockchain_vars.PEER_UPDATE_ACTION_ADD:
      magent_base_peer_updates.HandlePeerUpdateActionAdd(peerUpdate)
    case mblockchain_vars.PEER_UPDATE_ACTION_REMOVE:
      magent_base_peer_updates.HandlePeerUpdateActionRemove(peerUpdate)
    case mblockchain_vars.PEER_UPDATE_ACTION_IP_UPDATE:
      magent_base_peer_updates.HandlePeerUpdateActionIpUpdate(peerUpdate)
    default:
      mlog.GetLogger().Warnf("PeerUpdate with action %v not found", peerUpdate.Action)
    }
  }
}

/*
  Function loops indefinitely until a network contract address is filled in the config
*/
func (agent *AgentClient) idleForNetworkContractAddress() {
  const WAIT_SLEEP_TIME_S = 2
  for {
    networkContractAddress := mconfig.GetUserConfig().Blockchain.Network_Contract_Address
    if strings.Compare(networkContractAddress, mnet_vars.EMPTY_CONTRACT_ADDRESS) != 0 {
      break
    }
    mlog.GetLogger().Info("NetworkContractAddress is empty, waiting until there is a value")
    time.Sleep(time.Second * time.Duration(WAIT_SLEEP_TIME_S))
  }
}

/*
  Function loops indefinitely until a connection is established to the middleware process
*/
func (agent *AgentClient) idleForMiddlewareRegistration(selfPubKeyHash string) {
  // Register to middleware for peer updates, if the attempt fails, will attempt to try again
  const RETRY_SLEEP_TIME_S = 5
  for {
    mlog.GetLogger().Info("Attempting to subscribe to middleware")
    err := magent_middleware_interface.RegisterForPeerUpdates(selfPubKeyHash)
    if err == nil {
      mlog.GetLogger().Info("Subscribed to middleware")
      break
    }
    mlog.GetLogger().Warn("Could not connect to middleware, will re-attempt...")
    time.Sleep(time.Second * time.Duration(RETRY_SLEEP_TIME_S))
  }
}

/*
  Function starts a goroutine that blocks until a os.Signal is received in the teardownSignal channel.
  After the signal is received, do some additional cleanup for the agent
*/
func (agent *AgentClient) waitForTermSignal() {
  go func() {
    sig := <-agent.teardownSignal
    mlog.GetLogger().Info(fmt.Sprintf("Received os.Signal: %s", sig))
    mnet_core_manager.GetNetCoreManager().RemoveAllBridges()
    mlog.GetLogger().Info("Cleanup process completed. Exiting Marconid...")
    os.Exit(0)
  }()
}

/*
  Callback to handle when a peer response is received from the DHT
*/
func (agent *AgentClient) requestPeerResponseHandler(args map[string]string) {
  peerPubKeyHash := args["peerPubKeyHash"]
  host := args["peerIp"]

  if !mpeer.PeerManagerInstance().IsPeer(peerPubKeyHash) {
    // print warning if the received pkh is not our own
    if mcrypto_key.KeyManagerInstance().GetBasePublicKeyHash() != peerPubKeyHash {
      mlog.GetLogger().Warn(fmt.Sprintf("PeerRequestResponse received from non-peer: %s, noop", peerPubKeyHash))
    }
    return
  }
  mlog.GetLogger().Debug(fmt.Sprintf("PeerRequestResponse received for pubkeyhash: %s from %s", peerPubKeyHash, args["peerIp"]))

  // Check if the peer is already being handled
  agent.peerResponseHandlerStatusMutex.Lock()
  if attempting, exists := agent.peerResponseHandlerStatus[peerPubKeyHash]; !exists || (exists && !attempting) {
    agent.peerResponseHandlerStatus[peerPubKeyHash] = true
    agent.peerResponseHandlerStatusMutex.Unlock()

    // Note: requestPeerRouteBeacon can also return results from itself, ignore those
    if host != mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr() && host != "127.0.0.1" {

      logErrorAndUpdateStatus := func(err error) {
        mlog.GetLogger().Error(err)
        agent.peerResponseHandlerStatus[peerPubKeyHash] = false
      }

      // Initiate a public key exchange with the target peer (when necessary)
      err := mcrypto_key.KeyManagerInstance().InitiatePublicKeyExchange(host, peerPubKeyHash)
      if err != nil {
        logErrorAndUpdateStatus(err)
        return
      }

      // Initiate a Diffie-Hellman key exchange with the target peer (when necessary)
      err = mcrypto_dh.DHExchangeManagerInstance().InitiateDHKeyExchange(peerPubKeyHash)
      if err != nil {
        logErrorAndUpdateStatus(err)
        return
      }

      // Create a connection with the peer
      err = CreateConnectionToServicePeer(peerPubKeyHash, agent.baseL2KeyFilePath, true, host)
      if err != nil {
        logErrorAndUpdateStatus(err)
        return
      }

      // Stop the peer route request for this peerPubKeyHash
      mnet_dht.GetBeaconManager().StopPeerRouteRequest(peerPubKeyHash)
    }

    // Always clean up the status
    agent.peerResponseHandlerStatus[peerPubKeyHash] = false
  } else {
    agent.peerResponseHandlerStatusMutex.Unlock()
  }
}
