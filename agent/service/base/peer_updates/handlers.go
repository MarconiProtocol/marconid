package magent_base_peer_updates

import (
  "../../../../core/blockchain/vars"
  "../../../../core/crypto/dh"
  "../../../../core/crypto/key"
  "../../../../core/net/core/manager"
  "../../../../core/net/dht"
  "../../../../core/peer"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
)

/*
	Handle a PeerUpdate event of type mblockchain_vars.PEER_UPDATE_ACTION_ADD (ADD)
*/
func HandlePeerUpdateActionAdd(peerUpdate mblockchain_vars.PeerUpdate) {
  // Make sure the bridge for the main network is already created/allocated
  _, err := mnet_core_manager.GetNetCoreManager().GetBridgeInfoForNetwork(mnet_core_manager.SERVICE_NET, "main")
  if err != nil {
    mlog.GetLogger().Debug(fmt.Sprintf("Bridge not created, skipping: %s", err))
    return
  }

  // Safeguard for duplicated add peer action being received for the same PubKeyHash
  if mpeer.PeerManagerInstance().IsPeer(peerUpdate.PeerPubKeyHash) {
    mlog.GetLogger().Warn(fmt.Sprintf("Received PeerUpdate ADD for an existing peer with PubKeyHash %s", peerUpdate.PeerPubKeyHash))
    return
  }
  // add the peer
  mpeer.PeerManagerInstance().AddPeer(mpeer.SERVICE_PEER, peerUpdate.PeerPubKeyHash)
  // add the peer beacon in DHT
  mnet_dht.GetBeaconManager().StartPeerRouteRequest(peerUpdate.PeerPubKeyHash)
}

/*
	Handle a PeerUpdate event of type mblockchain_vars.PEER_UPDATE_ACTION_REMOVE (REMOVE)
*/
func HandlePeerUpdateActionRemove(peerUpdate mblockchain_vars.PeerUpdate) {
  dhManager := mcrypto_dh.DHExchangeManagerInstance()
  keyManager := mcrypto_key.KeyManagerInstance()

  peer, err := mpeer.PeerManagerInstance().GetPeer(peerUpdate.PeerPubKeyHash)
  if err != nil {
    mlog.GetLogger().Errorf("Failed to GetPeer %v, err = %v", peerUpdate.PeerPubKeyHash, err)
  } else {
    // remove the peer from peer manager
    mlog.GetLogger().Info("Removing peer from PeerManager")
    mpeer.PeerManagerInstance().RemovePeer(peerUpdate.PeerPubKeyHash)

    // RemoveAllBridges the Mpipe for the peer
    mnet_core_manager.GetNetCoreManager().CloseMPipe(peer)

    // stop the peer beacon in DHT
    mlog.GetLogger().Info("Calling StopPeerRouteRequest")
    mnet_dht.GetBeaconManager().StopPeerRouteRequest(peerUpdate.PeerPubKeyHash)
    // remove peer from DH exchange manager
    mlog.GetLogger().Info("Remove Peer from DHKeyInfo")
    dhManager.RemoveDHKeyInfo(peerUpdate.PeerPubKeyHash)
    // remove peer from key manager
    mlog.GetLogger().Info("Deleting Peer's public key")
    keyManager.DeletePeerPublicKey(peerUpdate.PeerPubKeyHash)

    mlog.GetLogger().Info("Removing peer with PubKeyHash completed", peer.PubKeyHash)
  }
}

/*
	Handle a PeerUpdate event of type mblockchain_vars.PEER_UPDATE_ACTION_IP_UPDATE (IP_UPDATE)
*/
func HandlePeerUpdateActionIpUpdate(peerUpdate mblockchain_vars.PeerUpdate) {
  nm := mnet_core_manager.GetNetCoreManager()
  // Check if the bridge has already been allocated, if not create one
  bridgeInfo, newlyAllocated, err := nm.GetOrAllocateBridgeInfoForNetwork(mnet_core_manager.SERVICE_NET, "main")
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("Could not get or allocate a bridge: %s", err))
  } else {
    if newlyAllocated {
      nm.CreateBridge(bridgeInfo, peerUpdate.IP, peerUpdate.NetMask, false)
    }
    nm.AssignIpAddrToBridge(bridgeInfo, peerUpdate.IP, peerUpdate.NetMask)
  }
}
