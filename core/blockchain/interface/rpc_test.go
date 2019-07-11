package mblockchain_interface

import (
  "../vars"
  "net/http"
  "reflect"
  "testing"
)

func TestBlockchainRPC_Init(t *testing.T) {
  tests := []struct {
    name  string
    bcRPC *BlockchainRPC
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.bcRPC.Init()
    })
  }
}

func TestBlockchainRPC_GetPeerUpdates(t *testing.T) {
  tests := []struct {
    name  string
    bcRPC *BlockchainRPC
    want  chan mblockchain_vars.PeerUpdate
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.bcRPC.GetPeerUpdates(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainRPC.GetPeerUpdates() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainRPC_GetEdgePeerUpdates(t *testing.T) {
  tests := []struct {
    name  string
    bcRPC *BlockchainRPC
    want  chan mblockchain_vars.EdgePeerUpdate
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.bcRPC.GetEdgePeerUpdates(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainRPC.GetEdgePeerUpdates() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainRPC_handleRpcUpdatePeers(t *testing.T) {
  type args struct {
    r           *http.Request
    w           http.ResponseWriter
    reqInfohash string
    reqPayload  string
  }
  tests := []struct {
    name  string
    bcRPC *BlockchainRPC
    args  args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.bcRPC.handleRpcUpdatePeers(tt.args.r, tt.args.w, tt.args.reqInfohash, tt.args.reqPayload)
    })
  }
}

func TestBlockchainRPC_handleRpcUpdateEdgePeers(t *testing.T) {
  type args struct {
    r           *http.Request
    w           http.ResponseWriter
    reqInfohash string
    reqPayload  string
  }
  tests := []struct {
    name  string
    bcRPC *BlockchainRPC
    args  args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.bcRPC.handleRpcUpdateEdgePeers(tt.args.r, tt.args.w, tt.args.reqInfohash, tt.args.reqPayload)
    })
  }
}

func Test_contains(t *testing.T) {
  type args struct {
    s []string
    e string
  }
  tests := []struct {
    name string
    args args
    want bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := contains(tt.args.s, tt.args.e); got != tt.want {
        t.Errorf("contains() = %v, want %v", got, tt.want)
      }
    })
  }
}
