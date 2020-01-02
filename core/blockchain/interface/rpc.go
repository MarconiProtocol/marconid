package mblockchain_interface

import (
  "../vars"
  mlog "github.com/MarconiProtocol/log"
  "net/http"
  "strconv"
  "strings"
)

type BlockchainRPC struct {
  peerUpdates     chan mblockchain_vars.PeerUpdate
  edgePeerUpdates chan mblockchain_vars.EdgePeerUpdate
  currentPeers    map[string]bool
}

type UpdatePeersArgs struct {
  PeersList string
  DHCP      string
  Active    bool
  ErrorCode bool
}

type UpdatePeersReply struct {
}

type UpdateEdgePeersArgs struct{}

type UpdateEdgePeersReply struct{}

func (bcRPC *BlockchainRPC) Init() {
  // instantiate peer update channels
  bcRPC.peerUpdates = make(chan mblockchain_vars.PeerUpdate)
  bcRPC.edgePeerUpdates = make(chan mblockchain_vars.EdgePeerUpdate)
  bcRPC.currentPeers = make(map[string]bool)
}

func (bcRPC *BlockchainRPC) GetPeerUpdates() chan mblockchain_vars.PeerUpdate {
  return bcRPC.peerUpdates
}

func (bcRPC *BlockchainRPC) GetEdgePeerUpdates() chan mblockchain_vars.EdgePeerUpdate {
  return bcRPC.edgePeerUpdates
}

/*
  Handler for incoming rpc requests to update peers
  Pushes specific PeerUpdates to the peerUpdates channel after parsing the rpc request
*/

func (bcRPC *BlockchainRPC) UpdatePeers(r *http.Request, args *UpdatePeersArgs, reply *UpdatePeersReply) error {
  // Assumption: PEERS is a comma-separated string such as abc,xyz,123. IP is a string on its own.
  // PeersList: 2be20ac5ce8c57ade93fdf3ee34e2ca6165dd551,e56fb56d8a91bf792e2f9951f25bdc2488a0fd9d
  // DHCP: 10.27.16.1/24
  if args.ErrorCode {
    mlog.GetLogger().Error("RPC UpdatePeers received error = true, which meant middleware's call to GetPeerInfo() caught an error")
    return nil
  }

  if !args.Active {
    // this peer no longer belong to the network, remove all current peers
    for peer := range bcRPC.currentPeers {
      peerUpdateRemove := mblockchain_vars.PeerUpdate{Action: mblockchain_vars.PEER_UPDATE_ACTION_REMOVE, PeerPubKeyHash: peer}
      bcRPC.peerUpdates <- peerUpdateRemove
      delete(bcRPC.currentPeers, peer)
    }
  } else {
    // peer still active handle its events
    peersList := strings.Split(args.PeersList, ",")
    dhcp := strings.Split(args.DHCP, "/")
    if len(dhcp) != 2 {
      mlog.GetLogger().Error("RPC UpdatePeers received payload with invalid DHCP", dhcp)
      return nil
    }
    ip := dhcp[0]
    netMask, err := strconv.Atoi(dhcp[1])
    if err != nil {
      mlog.GetLogger().Error("RPC UpdatePeers failed to retrieve net-mask", err)
      return nil
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
  return nil
}

/*
  TODO: STUB
  Handler for incoming rpc requests of type rpc.UPDATE_EDGE_PEERS ( "rpcUpdateEdgePeers" )
  Pushes specific EdgePeerUpdates to the edgePeerUpdates channel after parsing the rpc request
*/
func (bcRPC *BlockchainRPC) UpdateEdgePeers(r *http.Request, args *UpdateEdgePeersArgs, reply *UpdateEdgePeersReply) error {
  // STUB
  return nil
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
