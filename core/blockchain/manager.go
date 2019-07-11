// Package blockchain provides the functions to access and manage the Marconi blockchain and encapsulates
// any states which we build from data on the blockchain.
package mblockchain

import (
  "./interface"
  "./vars"
  "sync"
)

/*
  Blockchain data source interface
*/
type BlockchainDataSource interface {
  Init()
  GetPeerUpdates() chan mblockchain_vars.PeerUpdate
  GetEdgePeerUpdates() chan mblockchain_vars.EdgePeerUpdate
}

type BlockchainManager struct {
  networkId int
  nodeMap   map[string]mblockchain_vars.NodeInfo // a map of nodes where key is the node id in format of #n where n is a number, value is NodeInfo containing information such as ID and IP of the node
  meshMap   map[string][]string                  // a map of node and their respective peers, key is the node id, and value is a list of node ids that the node is connected to

  dataSource BlockchainDataSource
}

var blockchainManager *BlockchainManager
var once sync.Once

// returns the singleton to BlockchainManager
func GetBlockchainManager() *BlockchainManager {
  once.Do(func() {
    blockchainManager = &BlockchainManager{}
    blockchainManager.init()
  })
  return blockchainManager
}

/*
  Initialize the blockchain manager singleton instance
  The datasource created depends on configuration
*/
func (blockchainManager *BlockchainManager) init() {
  blockchainManager.dataSource = &mblockchain_interface.BlockchainRPC{}
  blockchainManager.dataSource.Init()
}

/*
  Return a channel that will produce PeerUpdates
*/
func (blockchainManager *BlockchainManager) GetPeerUpdates() chan mblockchain_vars.PeerUpdate {
  return blockchainManager.dataSource.GetPeerUpdates()
}

/*
  Return a channel that will produce EdgePeerUpdates
*/
func (BlockchainManager *BlockchainManager) GetEdgePeerUpdates() chan mblockchain_vars.EdgePeerUpdate {
  return blockchainManager.dataSource.GetEdgePeerUpdates()
}

// input is the node id
// return is a list of NodeInfo that the node is connected to
func (blockchainManager *BlockchainManager) GetNode(id string) mblockchain_vars.NodeInfo {
  return blockchainManager.nodeMap[id]
}

// input is the node id
// return is a list of node ids that the node is connected to
func (blockchainManager *BlockchainManager) GetMeshList(id string) []string {
  return blockchainManager.meshMap[id]
}
