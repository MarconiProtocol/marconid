package main

import (
  "./base"
  "flag"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
)

func main() {
  baseBeaconFilePath := flag.String("baseroutekey", "", "File path to the base route beacon key")
  L2keyFilePath := flag.String("l2key", "", "File path to the L2Key")
  baseDir := flag.String("basedir", "", "Base directory")
  edgeBeaconFilePath := flag.String("edgeroutekey", "", "Optional file path to the edge route beacon key")
  flag.Parse()

  mlog.GetLogger().Info(fmt.Sprintf("marconid invoked with args:"))
  mlog.GetLogger().Info(fmt.Sprintf("baseBeacon: %s", *baseBeaconFilePath))
  mlog.GetLogger().Info(fmt.Sprintf("L2KeyFilePath: %s", *L2keyFilePath))
  mlog.GetLogger().Info(fmt.Sprintf("baseDir: %s", *baseDir))
  mlog.GetLogger().Info(fmt.Sprintf("edgeBeaconFilePath: %s", *edgeBeaconFilePath))

  if *baseBeaconFilePath == "" {
    mlog.GetLogger().Fatal("Missing path to base beacon, cmd flag -baseroutekey")
  } else if *L2keyFilePath == "" {
    mlog.GetLogger().Fatal("Missing path to L2 key, cmd flag -l2key")
  } else if *baseDir == "" {
    mlog.GetLogger().Fatal("Missing base dir, cmd flag -basedir")
  }

  config := &magent_base.AgentConfig{
    BaseRouteBeaconKeyFilePath: *baseBeaconFilePath,
    L2KeyFilePath:              *L2keyFilePath,
    BaseDir:                    *baseDir,
    EdgeRouteBeaconKeyFilePath: *edgeBeaconFilePath,
  }
  //socksServer, err := socks5.Initialize()
  //if err != nil {
  //  mlog.GetLogger().Error("Failed to initialize the socks5 server in service::main!")
  //} else {
  //  socksServer.Start()
  //}
  agentClient := magent_base.NewAgentClient(config)
  agentClient.Start()
}
