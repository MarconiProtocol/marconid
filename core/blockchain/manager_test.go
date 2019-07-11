// Package blockchain provides the functions to access and manage the Marconi blockchain and encapsulates
// any states which we build from data on the blockchain.
package mblockchain

import (
  "./vars"
  "reflect"
  "testing"
)

func TestGetBlockchainManager(t *testing.T) {
  tests := []struct {
    name string
    want *BlockchainManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetBlockchainManager(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetBlockchainManager() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainManager_init(t *testing.T) {
  tests := []struct {
    name              string
    blockchainManager *BlockchainManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.blockchainManager.init()
    })
  }
}

func TestBlockchainManager_GetPeerUpdates(t *testing.T) {
  tests := []struct {
    name              string
    blockchainManager *BlockchainManager
    want              chan mblockchain_vars.PeerUpdate
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.blockchainManager.GetPeerUpdates(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainManager.GetPeerUpdates() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainManager_GetEdgePeerUpdates(t *testing.T) {
  tests := []struct {
    name              string
    BlockchainManager *BlockchainManager
    want              chan mblockchain_vars.EdgePeerUpdate
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.BlockchainManager.GetEdgePeerUpdates(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainManager.GetEdgePeerUpdates() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainManager_GetNode(t *testing.T) {
  type args struct {
    id string
  }
  tests := []struct {
    name              string
    blockchainManager *BlockchainManager
    args              args
    want              mblockchain_vars.NodeInfo
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.blockchainManager.GetNode(tt.args.id); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainManager.GetNode() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlockchainManager_GetMeshList(t *testing.T) {
  type args struct {
    id string
  }
  tests := []struct {
    name              string
    blockchainManager *BlockchainManager
    args              args
    want              []string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.blockchainManager.GetMeshList(tt.args.id); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BlockchainManager.GetMeshList() = %v, want %v", got, tt.want)
      }
    })
  }
}
