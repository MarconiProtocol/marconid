package mpeer

import (
  "strings"
  "testing"
)

var pm *PeerManager
var peer *Peer
var pubKeyHash = "0a2be4d884e10e4d0151"
var ip = "10.27.16.1"

func TestPeerManagerInstance(t *testing.T) {
  t.Log("TestPeerManagerInstance called")
  pm = PeerManagerInstance()
  if pm != PeerManagerInstance() {
    t.Error("Can only have one PeerManager instance")
  }
}

func TestPeerManager_AddPeer(t *testing.T) {
  t.Log("TestPeerManager_AddPeer called")
  peer = pm.AddPeer(SERVICE_PEER, pubKeyHash)
  if peer.Type != SERVICE_PEER || strings.Compare(peer.PubKeyHash, pubKeyHash) != 0 {
    t.Error("AddPeer failed")
  }

  if !pm.IsPeer(pubKeyHash) {
    t.Error("AddPeer failed")
  }
}

func TestPeerManager_IsPeer(t *testing.T) {
  t.Log("TestPeerManager_IsPeer called")
  if !pm.IsPeer(pubKeyHash) {
    t.Error("IsPeer failed")
  }
}

func TestPeerManager_GetPeer(t *testing.T) {
  t.Log("TestPeerManager_GetPeer called")
  peer1, err := pm.GetPeer(pubKeyHash)
  if err != nil || peer != peer1 {
    t.Error("GetPeer failed")
  }
}

func TestPeerManager_GetPeerOfType(t *testing.T) {
  t.Log("TestPeerManager_GetPeerOfType called")
  peer1, err := pm.GetPeerOfType(SERVICE_PEER, pubKeyHash)
  if peer != peer1 || err != nil {
    t.Error("GetPeerOfType failed")
  }
}

func TestPeerManager_UpdatePeer(t *testing.T) {
  t.Log("TestPeerManager_UpdatePeer called")
  pm.UpdatePeer(pubKeyHash, ip)
  if strings.Compare(peer.Ip, ip) != 0 {
    t.Error("UpdatePeer failed")
  }
}

func TestPeerManager_RemovePeer(t *testing.T) {
  t.Log("TestPeerManager_RemovePeer called")
  pm.RemovePeer(pubKeyHash)
  if pm.IsPeer(pubKeyHash) {
    t.Error("RemovePeer failed")
  }

  if _, err := pm.GetPeer(pubKeyHash); err == nil {
    t.Error("RemovePeer failed")
  }

  if _, err := pm.GetPeerOfType(SERVICE_PEER, pubKeyHash); err == nil {
    t.Error("RemovePeer failed")
  }
}
