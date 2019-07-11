package mblockchain_interface

import (
  "../../rpc"
  "../vars"
  "gitlab.neji.vm.tc/marconi/log"
  "net/http"
  "strconv"
  "strings"
)

type BlockchainRPC struct {
  incomingRPCPayloads chan string

  peerUpdates     chan mblockchain_vars.PeerUpdate
  edgePeerUpdates chan mblockchain_vars.EdgePeerUpdate
  currentPeers    map[string]bool
}

func (bcRPC *BlockchainRPC) Init() {
  // instantiate peer update channels
  bcRPC.peerUpdates = make(chan mblockchain_vars.PeerUpdate)
  bcRPC.edgePeerUpdates = make(chan mblockchain_vars.EdgePeerUpdate)

  bcRPC.currentPeers = make(map[string]bool)

  // Register RPC handlers
  rpc.RegisterRpcHandler(rpc.UPDATE_PEERS, bcRPC.handleRpcUpdatePeers)
  rpc.RegisterRpcHandler(rpc.UPDATE_EDGE_PEERS, bcRPC.handleRpcUpdateEdgePeers)
}

func (bcRPC *BlockchainRPC) GetPeerUpdates() chan mblockchain_vars.PeerUpdate {
  return bcRPC.peerUpdates
}

func (bcRPC *BlockchainRPC) GetEdgePeerUpdates() chan mblockchain_vars.EdgePeerUpdate {
  return bcRPC.edgePeerUpdates
}

/*
  Handler for incoming rpc requests of type rpc.UPDATE_PEERS ( "rpcUpdatePeers" )
  Pushes specific PeerUpdates to the peerUpdates channel after parsing the rpc request
*/
func (bcRPC *BlockchainRPC) handleRpcUpdatePeers(r *http.Request, w http.ResponseWriter, reqInfohash, reqPayload string) {
  // Assumption: reqPayload is in the format of PEERS;IP, where PEERS is a comma-separated string
  // such as abc,xyz,123. IP is a string on its own. This can be improved later when middleware pass payload in JSON.
  // 2be20ac5ce8c57ade93fdf3ee34e2ca6165dd551,e56fb56d8a91bf792e2f9951f25bdc2488a0fd9d;10.27.16.1/24;true
  info := strings.Split(reqPayload, ";")

  if len(info) != 3 {
    mlog.GetLogger().Error("RPC handleRpcUpdatePeers received invalid payload", reqPayload)
    return
  }

  active, err := strconv.ParseBool(info[2])
  if err != nil {
    mlog.GetLogger().Error("RPC handleRpcUpdatePeers failed to retrieve active status", err)
    return
  }
  if !active {
    // this peer no longer belong to the network, remove all current peers
    for peer := range bcRPC.currentPeers {
      peerUpdateRemove := mblockchain_vars.PeerUpdate{Action: mblockchain_vars.PEER_UPDATE_ACTION_REMOVE, PeerPubKeyHash: peer}
      bcRPC.peerUpdates <- peerUpdateRemove
      delete(bcRPC.currentPeers, peer)
    }
  } else {
    // peer still active handle its events
    peersList := strings.Split(info[0], ",")
    dhcp := strings.Split(info[1], "/")
    if len(dhcp) != 2 {
      mlog.GetLogger().Error("RPC handleRpcUpdatePeers received payload with invalid DHCP", dhcp)
      return
    }
    ip := dhcp[0]
    netMask, err := strconv.Atoi(dhcp[1])
    if err != nil {
      mlog.GetLogger().Error("RPC handleRpcUpdatePeers failed to retrieve net-mask", err)
      return
    }

    // handle peer IP updates
    peerUpdateIP := mblockchain_vars.PeerUpdate{Action: mblockchain_vars.PEER_UPDATE_ACTION_IP_UPDATE, IP: ip, NetMask: netMask}
    bcRPC.peerUpdates <- peerUpdateIP

    // handle removed peers
    for peer := range bcRPC.currentPeers {
      if !contains(peersList, peer) {
        peerUpdateRemove := mblockchain_vars.PeerUpdate{Action: mblockchain_vars.PEER_UPDATE_ACTION_REMOVE, PeerPubKeyHash: peer}
        bcRPC.peerUpdates <- peerUpdateRemove
        delete(bcRPC.currentPeers, peer)
      }
    }

    // handle new peers
    for _, peer := range peersList {
      if _, exists := bcRPC.currentPeers[peer]; !exists {
        peerUpdateAdd := mblockchain_vars.PeerUpdate{Action: mblockchain_vars.PEER_UPDATE_ACTION_ADD, PeerPubKeyHash: peer}
        bcRPC.peerUpdates <- peerUpdateAdd
        bcRPC.currentPeers[peer] = true
      }
    }
  }
}

/*
  TODO: STUB
  Handler for incoming rpc requests of type rpc.UPDATE_EDGE_PEERS ( "rpcUpdateEdgePeers" )
  Pushes specific EdgePeerUpdates to the edgePeerUpdates channel after parsing the rpc request
*/
func (bcRPC *BlockchainRPC) handleRpcUpdateEdgePeers(r *http.Request, w http.ResponseWriter, reqInfohash, reqPayload string) {
  // STUB
}

// check if a slice contains an element
func contains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}
