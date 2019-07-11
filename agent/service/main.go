package main

import (
  "./base"
  "fmt"
  "os"
)

func main() {
  fmt.Println("marconid invoked with args:")
  for i, arg := range os.Args {
    fmt.Println(i, ": ", arg)
  }
  config := &magent_base.AgentConfig{
    os.Args[2],
    os.Args[1],
  }
  agentClient := magent_base.NewAgentClient(config)
  agentClient.Start()
}
