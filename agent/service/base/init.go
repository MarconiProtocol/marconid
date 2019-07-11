package magent_base

import (
  "../../../core"
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/ip"
  "../../../core/net/packet/filter"
  "../../../core/rpc"
  "../../../core/runtime"
  "../../../util"
  "fmt"
  mlog "gitlab.neji.vm.tc/marconi/log"
  "os"
  "os/signal"
  "syscall"
)

/*
  Configuration object for a Marconi agent
*/
type AgentConfig struct {
  BaseRouteBeaconKeyFilePath string
  L2KeyFilePath              string
}

/*
  Initialize an agent client object with an agent config
*/
func (agent *AgentClient) initialize(conf *AgentConfig) {
  mcore.Initialize()
  agent.initInterrupt()

  agent.peerResponseHandlerStatus = make(map[string]bool)

  // Initialize runtime object defaults
  mrt := mruntime.GetMRuntime()
  mrt.SetRuntimeOS()
  mainIpAddr, err := mnet_ip.GetMainInterfaceIpAddress()
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("Failed getting the main network interface ip address: %s", err))
  }
  mrt.InterfaceInfo.SetLocalMainInterfaceIpAddr(mainIpAddr)
  agent.baseL2KeyFilePath = conf.L2KeyFilePath
  agent.baseBeaconKey = mutil.LoadKey(conf.BaseRouteBeaconKeyFilePath)

  // Effectively initialize the managers
  mpacket_filter.GetFilterManagerInstance()
  mcrypto_dh.DHExchangeManagerInstance()

  // Ensure that private/public keys were generated
  mcrypto_key.KeyManagerInstance().EnsurePrivatePublicKeysGenerated()
  mcrypto_key.KeyManagerInstance().LoadDefaultBaseKey()

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
