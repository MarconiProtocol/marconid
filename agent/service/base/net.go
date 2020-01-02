package magent_base

import (
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/core/manager"
  "../../../core/net/vars"
  "../../../core/runtime"
  "../../../util"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "net"
  "strconv"
  "strings"
)

/*
  Attempt to create a connection to a service node peer.
  This function is invoked once the public key and dh key exchanges are completed as a callback
*/
func CreateConnectionToServicePeer(peerPubKeyHash string, keyFile string, isSecure bool, peerIp string) error {
  keyManager := mcrypto_key.KeyManagerInstance()
  localTargetPort := mutil.GetMutualMPipePort(keyManager.GetBasePublicKeyHash(), peerPubKeyHash)

  // Check to see if the peer is a valid potential peer
  peerIp = strings.TrimSpace(peerIp)

  // Get DHKeyInfo for the peer
  dhKeyInfo, err := mcrypto_dh.DHExchangeManagerInstance().GetDHKeyInfo(peerPubKeyHash)
  if err != nil {
    return errors.New(fmt.Sprintf("Could not find DHKeyInfo for peer: %s ", peerPubKeyHash))
  }

  // Checks if the port is in use
  if exists, connID := mnet_core_manager.GetNetCoreManager().CheckConnectionPortExists(strconv.Itoa(localTargetPort)); exists {
    mlog.GetLogger().Info(fmt.Sprintf("Port %d is already used by a connection %d", localTargetPort, connID))
    return nil
  }

  // Try to allocate a udp port (localTargetPort)
  addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localTargetPort))
  ln, err := net.ListenUDP("udp", addr)
  if err == nil {
    err = ln.Close()
    if err == nil {

      mlog.GetLogger().Debug(fmt.Sprintf("CreateConnectionToServicePeer: CreateMPipe: %d", localTargetPort))
      mlog.GetLogger().Debug(fmt.Sprintf(" - Remote Peer IP: %s \n", peerIp))
      mlog.GetLogger().Debug(fmt.Sprintf(" - Local Main interface IP: %s \n", mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr()))

      err := mnet_core_manager.GetNetCoreManager().CreateMPipe(
        &mnet_vars.ConnectionArgs{
          L2KeyFile:    keyFile,
          EncPayload:   isSecure,
          LocalPort:    strconv.Itoa(localTargetPort),
          RemoteIpAddr: peerIp,
          RemotePort:   strconv.Itoa(localTargetPort), // We use the same remote port as the local port because each side will deterministically get the port number
          // this however means the port must be open, otherwise it is an error case
          DataKey:        dhKeyInfo.SymmetricKeyBytes,
          DataKeySignal:  dhKeyInfo.KeySignal,
          PeerPubKeyHash: peerPubKeyHash,
        })
      if err != nil {
        return errors.New(fmt.Sprintf("Error calling CreateMPipe: %s", err))
      }
    }
  } else {
    return errors.New(fmt.Sprintf("Port %d is used by another process. Error: %s", localTargetPort, err))
  }
  return nil
}
