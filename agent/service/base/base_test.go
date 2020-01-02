package magent_base

import (
  "reflect"
  "testing"
)

func TestNewAgentClient(t *testing.T) {
  type args struct {
    conf *AgentConfig
  }
  tests := []struct {
    name string
    args args
    want *AgentClient
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := NewAgentClient(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("NewAgentClient() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestAgentClient_Start(t *testing.T) {
  tests := []struct {
    name  string
    agent *AgentClient
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.agent.Start()
    })
  }
}

func TestAgentClient_idleForNetworkContractAddress(t *testing.T) {
  tests := []struct {
    name  string
    agent *AgentClient
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.agent.idleForNetworkContractAddress()
    })
  }
}

func TestAgentClient_idleForMiddlewareRegistration(t *testing.T) {
  type args struct {
    selfPubKeyHash string
  }
  tests := []struct {
    name  string
    agent *AgentClient
    args  args
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.agent.idleForMiddlewareRegistration(tt.args.selfPubKeyHash)
    })
  }
}

func TestAgentClient_waitForTermSignal(t *testing.T) {
  tests := []struct {
    name  string
    agent *AgentClient
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.agent.waitForTermSignal()
    })
  }
}

func TestAgentClient_requestPeerResponseHandler(t *testing.T) {
  type args struct {
    args map[string]string
  }
  tests := []struct {
    name  string
    agent *AgentClient
    args  args
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.agent.requestPeerResponseHandler(tt.args.args)
    })
  }
}
