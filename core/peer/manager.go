package mpeer

import (
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "sync"
)

// TODO: Merge management of peer pub keys and dh key info here

type PeerType int

const (
  SERVICE_PEER PeerType = iota
  EDGE_PEER
)

type PeerManager struct {
  peers map[string]*Peer
}

type Peer struct {
  Type       PeerType
  PubKeyHash string
  Ip         string
  Port       string
}

var instance *PeerManager
var once sync.Once

func PeerManagerInstance() *PeerManager {
  once.Do(func() {
    instance = &PeerManager{}
    instance.initialize()
  })
  return instance
}

func (pm *PeerManager) initialize() {
  pm.peers = make(map[string]*Peer)
}

func (pm *PeerManager) AddPeer(peerType PeerType, peerPubKeyHash string) *Peer {
  if !pm.IsPeer(peerPubKeyHash) {
    newPeer := Peer{}
    newPeer.Type = peerType
    newPeer.PubKeyHash = peerPubKeyHash
    pm.peers[peerPubKeyHash] = &newPeer
  }
  return pm.peers[peerPubKeyHash]
}

func (pm *PeerManager) RemovePeer(peerPubKeyHash string) {
  if pm.IsPeer(peerPubKeyHash) {
    delete(pm.peers, peerPubKeyHash)
  }
}

func (pm *PeerManager) GetPeer(peerPubKeyHash string) (*Peer, error) {
  if !pm.IsPeer(peerPubKeyHash) {
    return nil, errors.New(fmt.Sprintf("No peer found with pubKeyHash %s", peerPubKeyHash))
  }
  return pm.peers[peerPubKeyHash], nil
}

func (pm *PeerManager) GetPeerOfType(peerType PeerType, peerPubKeyHash string) (*Peer, error) {
  peer, err := pm.GetPeer(peerPubKeyHash)
  if err != nil {
    return nil, err
  }
  if peer.Type != peerType {
    return nil, errors.New(fmt.Sprintf("Peer found with pubKeyHash %s found, but not of type %d", peerPubKeyHash, peerType))
  }
  return peer, nil
}

func (pm *PeerManager) UpdatePeer(peerPubKeyHash string, ip string) {
  if pm.IsPeer(peerPubKeyHash) {
    pm.peers[peerPubKeyHash].Ip = ip
  }
}

func (pm *PeerManager) IsPeer(peerPubKeyHash string) bool {
  _, exists := pm.peers[peerPubKeyHash]
  return exists
}

func (pm *PeerManager) DebugPrint() {
  mlog.GetLogger().Debug("PEERMANAGER - CURRENT PEERS")
  for pubKeyHash, peer := range pm.peers {
    mlog.GetLogger().Debug("Pubkeyhash", pubKeyHash, " peerObj", peer)
  }
}
