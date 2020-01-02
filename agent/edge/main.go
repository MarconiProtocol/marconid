package main

import (
  "./base"
  "flag"
)

func main() {
  l2KeyFilePath := flag.String("l2key", "", "The full path to the l2 key file")
  edgeRouteKeyFilePath := flag.String("edgeroutekey", "", "The full path to the edge route key file")
  id := flag.String("id", "", "Some kind of identifier for the FeO - TMP")
  baseDir := flag.String("basedir", "/opt/marconi", "The base directory")
  flag.Parse()

  // Create edge client config based on user input
  config := &magent_edge_base.AgentEdgeConfig{
    EdgeRouteBeaconKeyFilePath: *edgeRouteKeyFilePath,
    L2KeyFilePath:              *l2KeyFilePath,
    NetworkId:                  *id,
    BaseDir:                    *baseDir,
  }

  // Create and start a new edge client
  agentEdgeClient := magent_edge_base.NewAgentEdgeClient(config)
  agentEdgeClient.Start()
}
