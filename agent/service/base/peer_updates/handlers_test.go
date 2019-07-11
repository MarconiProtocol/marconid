package magent_base_peer_updates

import (
  "../../../../core/blockchain/vars"
  "testing"
)

func TestHandlePeerUpdateActionAdd(t *testing.T) {
  type args struct {
    peerUpdate mblockchain_vars.PeerUpdate
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      HandlePeerUpdateActionAdd(tt.args.peerUpdate)
    })
  }
}

func TestHandlePeerUpdateActionRemove(t *testing.T) {
  type args struct {
    peerUpdate mblockchain_vars.PeerUpdate
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      HandlePeerUpdateActionRemove(tt.args.peerUpdate)
    })
  }
}

func TestHandlePeerUpdateActionIpUpdate(t *testing.T) {
  type args struct {
    peerUpdate mblockchain_vars.PeerUpdate
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      HandlePeerUpdateActionIpUpdate(tt.args.peerUpdate)
    })
  }
}
