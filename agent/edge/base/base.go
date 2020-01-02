package magent_edge_base

import (
  "../../../core/net/dht"
  "../../../core/runtime"
  "../../../util"
  "../edge_route"
  "../node_manager"
  "crypto/rsa"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "os"
  "time"
)

const NUMBER_REQUESTED_TOP_NODES = 3
const NODE_SEARCH_INTERVAL_MS = 2000
const NODE_RETRY_INTERVAL_MS = 3000

type AgentEdgeClient struct {
  teardownSignal    chan os.Signal
  edgeBeaconKey     *rsa.PrivateKey
  baseL2KeyFilePath string
  networkId         string
}

/*
Construct the agent edge client, given a config
*/
func NewAgentEdgeClient(conf *AgentEdgeConfig) *AgentEdgeClient {
  agentEdgeClient := AgentEdgeClient{}
  agentEdgeClient.initialize(conf)
  return &agentEdgeClient
}

/*
  Start the edge client agent and establish an edge route
*/
func (agent *AgentEdgeClient) Start() {
  // Get the edge edgePKH
  edgePKH, _ := mutil.GetInfohashByPubKey(&agent.edgeBeaconKey.PublicKey)

  // Create an edge route beacon, and request for an edge peer
  mnet_dht.GetBeaconManager().CreateEdgeRouteBeacon(agent.requestEdgeRouteResponseHandler)
  mnet_dht.GetBeaconManager().StartEdgeRouteRequest(&agent.edgeBeaconKey.PublicKey)

  // Request edge route from edge peer(s) until one connected
  for {
    err := agent.EstablishEdgeRoute(edgePKH)
    if err == nil {
      mlog.GetLogger().Info("AgentEdgeClient::Start - Edge route established. Stopping Edge Route Requests.")
      break
    }
    mlog.GetLogger().Warn("Failed to connect to any of top nodes. Retrying.")
    time.Sleep(time.Millisecond * time.Duration(NODE_RETRY_INTERVAL_MS))
  }

  mnet_dht.GetBeaconManager().StopEdgeRouteRequest()

  // TODO: this should be updated to loop on commands from UI / or maybe other input.. we will see. Maybe re-request for an edge route on disconnect
  // For now the main goroutine will idle, while the mpipe goroutines execute
  select {}
}

/*
  Search for edge peer(s) and establish a route. Will only create one connection.
*/
func (agent *AgentEdgeClient) EstablishEdgeRoute(edgePKH string) error {

  // Start waiting until at least N nodes are found, then returns the top N nodes
  // sleeps for searchIntervalMS between checks of top nodes
  topPeers := magent_edge_node_manager.Instance().WaitForTopNodes(NUMBER_REQUESTED_TOP_NODES, NODE_SEARCH_INTERVAL_MS)

  // Make a request to create an edge route to a peer
  for _, peer := range topPeers {
    err := magent_edge_route_request.Instance().RequestRoute(peer, edgePKH)
    if err == nil {
      // TODO: ayuen - assumption for now, we want only one connection, after we create an edge route we are done
      return nil
    }
  }
  err := fmt.Errorf("failed to connect to any of %d top nodes. Connection not established", NUMBER_REQUESTED_TOP_NODES)
  return err
}

/*
  Callback to handle when a base route response is received from the DHT
*/
func (agent *AgentEdgeClient) requestEdgeRouteResponseHandler(args map[string]string) {
  pkh, _ := mutil.GetInfohashByPubKey(&agent.edgeBeaconKey.PublicKey)
  if args["peerIp"] != mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr() {
    magent_edge_node_manager.Instance().AddNode(args["peerIp"], pkh)
  }
}
