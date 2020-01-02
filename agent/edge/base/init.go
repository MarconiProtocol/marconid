package magent_edge_base

import (
  "../../../core"
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/core/manager"
  "../../../core/net/ip"
  "../../../core/rpc"
  "../../../core/runtime"
  "../../../util"
  "fmt"
  "github.com/MarconiProtocol/log"
  "os"
  "os/signal"
  "syscall"
)

type AgentEdgeConfig struct {
  BaseRouteBeaconKeyFilePath string
  EdgeRouteBeaconKeyFilePath string
  L2KeyFilePath              string
  NetworkId                  string
  BaseDir                    string
}

// ------
// TODO: ayuen cleanup
func (agent *AgentEdgeClient) initInterrupt() {
  agent.teardownSignal = make(chan os.Signal)
  signal.Notify(agent.teardownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
  agent.waitForTermSignal()
}

func (agent *AgentEdgeClient) waitForTermSignal() {
  go func() {
    sig := <-agent.teardownSignal
    mlog.GetLogger().Info(fmt.Sprintf("Received os.Signal: %s", sig))
    mnet_core_manager.GetNetCoreManager().RemoveAllBridges()
    mlog.GetLogger().Info("Cleanup process completed. Exiting Marconid...")
    os.Exit(0)
  }()
}

// ------

func (agent *AgentEdgeClient) initialize(conf *AgentEdgeConfig) {
  mcore.Initialize(conf.BaseDir)
  agent.initInterrupt()

  agent.networkId = conf.NetworkId

  //agent.peerResponseHandlerStatus = make(map[string]bool)

  // Initialize runtime object defaults
  mrt := mruntime.GetMRuntime()
  mrt.SetRuntimeOS()
  mainIpAddr, err := mnet_ip.GetMainInterfaceIpAddress()
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("Failed getting the main network interface ip address: %s", err))
  }
  mrt.InterfaceInfo.SetLocalMainInterfaceIpAddr(mainIpAddr)
  agent.baseL2KeyFilePath = conf.L2KeyFilePath
  agent.edgeBeaconKey = mutil.LoadKey(conf.EdgeRouteBeaconKeyFilePath)

  // Effectively initialize the managers
  //mpacket_filter.GetFilterManagerInstance()
  mcrypto_dh.DHExchangeManagerInstance()

  // Ensure that private/public keys were generated
  mcrypto_key.KeyManagerInstance().EnsurePrivatePublicKeysGenerated()
  mcrypto_key.KeyManagerInstance().LoadDefaultBaseKey()

  // Start the RPC server
  go rpc.StartRPCServer()
}
