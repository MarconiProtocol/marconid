package magent_base

import (
  "../../../core"
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/ip"
  "../../../core/net/packet/filter"
  "../../../core/net/tc"
  "../../../core/rpc"
  "../../../core/runtime"
  "../../../util"
  magent_util "../../service/util"
  "../edge_route"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "os"
  "os/signal"
  "syscall"
)

/*
  Configuration object for a Marconi agent
*/
type AgentConfig struct {
  BaseRouteBeaconKeyFilePath string
  EdgeRouteBeaconKeyFilePath string
  L2KeyFilePath              string
  BaseDir                    string
}

var Agent *AgentClient

/*
  Initialize an agent client object with an agent config
*/
func (agent *AgentClient) initialize(conf *AgentConfig) {
  mcore.Initialize(conf.BaseDir)
  agent.initInterrupt()

  agent.peerResponseHandlerStatus = make(map[string]bool)

  // Initialize runtime object defaults
  mrt := mruntime.GetMRuntime()
  mrt.SetRuntimeOS()
  mainIpAddr, err := mnet_ip.GetMainInterfaceIpAddress()
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("Failed getting the main network interface ip address: %s", err.Error()))
  }
  mrt.InterfaceInfo.SetLocalMainInterfaceIpAddr(mainIpAddr)
  agent.baseL2KeyFilePath = conf.L2KeyFilePath
  agent.baseBeaconKey = mutil.LoadKey(conf.BaseRouteBeaconKeyFilePath)
  if conf.EdgeRouteBeaconKeyFilePath != "" {
    agent.edgeBeaconKey = mutil.LoadKey(conf.EdgeRouteBeaconKeyFilePath)
  }

  // Effectively initialize the managers
  mpacket_filter.GetFilterManagerInstance()
  mcrypto_dh.DHExchangeManagerInstance()

  // Ensure that private/public keys were generated
  mcrypto_key.KeyManagerInstance().EnsurePrivatePublicKeysGenerated()
  mcrypto_key.KeyManagerInstance().LoadDefaultBaseKey()

  // Register the edge route request rpc handler
  edgeConnectionService := new(magent_edge_route_provider.EdgeConnectionService)
  if err := rpc.RegisterService(edgeConnectionService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for EdgeRouteManager")
  }

  // Register the node finder check rpc handler
  agentClientService := new(AgentClientService)
  if err := rpc.RegisterService(agentClientService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for AgentClient, err =", err)
  }

  // Register the traffic control rpc service
  trafficControlService := new(tc.TrafficControlService)
  if err := rpc.RegisterService(trafficControlService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for Traffic Control, err =", err)
  }

  // Register the netflow traffic monitor
  netflowService := new(magent_util.NetflowService)
  if err := rpc.RegisterService(netflowService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for Netflow Traffic Monitor, err =", err)
  }

  // Start the RPC server
  go rpc.StartRPCServer()
}

/*
  Start a goroutine to wait for a result on the teardown signal channel, which is waiting for sigint, sigterm or sigkill
*/
func (agent *AgentClient) initInterrupt() {
  agent.teardownSignal = make(chan os.Signal)
  signal.Notify(agent.teardownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
  agent.waitForTermSignal()
}
