package mblockchain_vars

import "net"

const (
  // peer update actions
  PEER_UPDATE_ACTION_ADD       = "ADD"
  PEER_UPDATE_ACTION_UPDATE    = "UPDATE"
  PEER_UPDATE_ACTION_REMOVE    = "REMOVE"
  PEER_UPDATE_ACTION_IP_UPDATE = "IP_UPDATE"

  // edge update actions
  EDGE_PEER_UPDATE_ACTION_SUBSCRIBE = "SUBSCRIBE"
  EDGE_PEER_UPDATE_ACTION_EXPIRE    = "EXPIRE"
  EDGE_PEER_UPDATE_ACTION_CANCEL    = "CANCEL"
)

type PeerUpdate struct {
  Action         string // Add for now, but we can add things like Remove/Update
  PeerPubKeyHash string
  IP             string
  NetMask        int
}

type EdgePeerUpdate struct {
  Action         string
  PeerPubKeyHash string
  Expiration     uint
}

type NodeInfo struct {
  ID        string // a unique identifier for the node, format is "#n" where n is a number
  IP        net.IP // ip of the node
  Hostname  string // hostname of the node
  TimeAdded string // timestamp of which the node was created
}
