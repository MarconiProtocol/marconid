package mcrypto_key

import (
  "../../peer"
  "../../rpc"
  "../../../util"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "encoding/base64"
  "errors"
  "gitlab.neji.vm.tc/marconi/log"
  "net/http"
  "strings"

  "fmt"
)


func (km *KeyManager) registerRPCHandlers() {
  rpc.RegisterRpcHandler(rpc.REQUEST_PUB_KEY_EXCHANGE_SYN, km.handlePubKeyExchangeSynRPC)
}

func (km *KeyManager) InitiatePublicKeyExchange(peerIp string, peerPubKeyHash string) error {
  // Try to get the Peer's key, if we don't have it, then we need to initiate public key exchange
  _, err := KeyManagerInstance().GetPeerPublicKey(peerPubKeyHash)
  if err != nil {
    mlog.GetLogger().Debug(fmt.Sprintf("PubKeyExchange::InitiatePublicKeyExchange - Host %s, PeerPubKeyHash %s", peerIp, peerPubKeyHash))
    // Start the PubKey Exchange process
    // Start the Syn assuming the target is not behind NAT
    return KeyManagerInstance().SendPubKeyExchangeSynRPC(peerIp, peerPubKeyHash)
  }

  return nil
}

func (km *KeyManager) SendPubKeyExchangeSynRPC(targetHost string, peerPubKeyHash string) error {
  mlog.GetLogger().Debug(fmt.Sprintf("SendPubKeyExchangeSynRPC - Sending syn to %s", targetHost))
  payload := buildPubKeyExchangePayload(peerPubKeyHash, km.GetBasePublicKey(), km.GetBasePrivateKey())

  response := rpc.SendRPC(targetHost, rpc.RPC_PORT, rpc.REQUEST_PUB_KEY_EXCHANGE_SYN, payload)
  if response.Error != "" {
    return errors.New(fmt.Sprintf("SendPubKeyExchangeSynRPC - Nil/error response from %s ; response %v", targetHost, response.Error))
  }

  mlog.GetLogger().Debug(fmt.Sprintf("SendPubKeyExchangeSynRPC - Received response from %s", targetHost))
  signature, pubKeyHash, peerPubKeyHash, encodedPeerPubKey := parsePubKeyExchangePayload(response.Result)
  if !km.isRpcSenderReceiverValid(peerPubKeyHash, pubKeyHash) {
    return errors.New(fmt.Sprintf("handlePubKeyExchangeAckRPC - invalid rpc sender, not a peer %s", peerPubKeyHash))
  }

  // Decode peerPubKey and verify signature
  peerPubKey := decodePublicKeyB64String(encodedPeerPubKey)
  err := verifySignature(peerPubKey, signature, peerPubKeyHash, pubKeyHash, encodedPeerPubKey)
  if err != nil {
    return errors.New(fmt.Sprintf("Signature verification failed on pubKeyExchangeSyn response."))
  }

  mpeer.PeerManagerInstance().UpdatePeer(peerPubKeyHash, targetHost)

  // Store public key of peer
  km.SetPeerPublicKey(peerPubKeyHash, peerPubKey)

  return nil
}

func (km *KeyManager) handlePubKeyExchangeSynRPC(r *http.Request, w http.ResponseWriter, reqInfohash string, reqPayload string) {
  // Using the sender IP in the http request; sender could be behind NAT
  ip, _, err := rpc.ParseIPAndPort(r.RemoteAddr)
  if err != nil {
    mlog.GetLogger().Error("PubKeyExchange::handlePubKeyExchangeSynRPC - could not parse remote address / port from request")
    return
  }

  signature, pubKeyHash, peerPubKeyHash, encodedPeerPubKey := parsePubKeyExchangePayload(reqPayload)
  if !km.isRpcSenderReceiverValid(peerPubKeyHash, pubKeyHash) {
    mlog.GetLogger().Debug("PubKeyExchange::handlePubKeyExchangeSynRPC - invalid rpc sender, sender is not a peer")
    err = rpc.WriteResponseWithError(&w, "Not recognized as a peer")
    return
  }
  mlog.GetLogger().Debug("PubKeyExchange::handlePubKeyExchangeSynRPC - RECEIVED SYN from", r.RemoteAddr)

  // Decode peerPubKey and verify signature
  peerPubKey := decodePublicKeyB64String(encodedPeerPubKey)
  err = verifySignature(peerPubKey, signature, peerPubKeyHash, pubKeyHash, encodedPeerPubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Signature verification failed.")
  }

  mpeer.PeerManagerInstance().UpdatePeer(peerPubKeyHash, ip)

  // Store public key of peer
  km.SetPeerPublicKey(peerPubKeyHash, peerPubKey)

  mlog.GetLogger().Debug("PubKeyExchange::SendPubKeyExchangeAckRPC - SENDING ACK to", ip)
  payload := buildPubKeyExchangePayload(peerPubKeyHash, km.GetBasePublicKey(), km.GetBasePrivateKey())
  err = rpc.WriteResponseWithResult(&w, payload)
  if err != nil {
    mlog.GetLogger().Debug(fmt.Sprintf("PubKeyExchange:: error writing to response: %s", err))
  }
}

func buildPubKeyExchangePayload(peerPubKeyHash string, pubKey *rsa.PublicKey, privateKey *rsa.PrivateKey) string {
  // craft our payload
  pubKeyHash, err := mutil.GetInfohashByPubKey(pubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to get pubkeyhash from pub key")
  }
  encodedPubKey := encodePubKeyToB64String(pubKey)
  signature := createSignature(privateKey, peerPubKeyHash, pubKeyHash, encodedPubKey)
  payload := strings.Join([]string{signature, peerPubKeyHash, pubKeyHash, encodedPubKey}, ";")
  // base64 encode the payload
  payload = base64.StdEncoding.EncodeToString([]byte(payload))
  return payload
}

func parsePubKeyExchangePayload(payload string) (signature string, pubKeyHash string, peerPubKeyHash string, encodedPubKey string) {
  // base64 decode the payload
  payloadBytes, err := base64.StdEncoding.DecodeString(payload)
  if err != nil {
    mlog.GetLogger().Error("Failed to decode pk exchange payload")
    return "", "", "", ""
  }
  payload = string(payloadBytes)
  // split the string payload into the signature, pubKeyHash, peerPubKeyHash, encodedPubKey
  // Note: peerPubKeyHash and pubKeyHash order is swapped on parse and build, simply because of perspective
  payloadSplit := strings.Split(payload, ";")
  return payloadSplit[0], payloadSplit[1], payloadSplit[2], payloadSplit[3]
}

func (km *KeyManager) isRpcSenderReceiverValid(peerPubKeyHash string, pubKeyHash string) bool {
  if !mpeer.PeerManagerInstance().IsPeer(peerPubKeyHash) {
    return false
  }

  pubKey := km.GetBasePublicKey()
  actualPubKeyHash, err := mutil.GetInfohashByPubKey(pubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to get pubkeyhash from pub key")
  }
  // Check if this RPC call is actually meant for us
  return pubKeyHash == actualPubKeyHash
}

// Sign our payload
func createSignature(privateKey *rsa.PrivateKey, peerPubKeyHash string, pubKeyHash string, encodedPubKey string) (signature string) {
  payloadJoined := strings.Join([]string{peerPubKeyHash, pubKeyHash, encodedPubKey}, ";")

  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  pssh := newhash.New()
  pssh.Write([]byte(payloadJoined))
  hashed := pssh.Sum(nil)

  signatureBytes, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, &opts)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to generate signature")
  }
  return base64.StdEncoding.EncodeToString(signatureBytes)
}

func verifySignature(peerPubKey *rsa.PublicKey, signatureEncoded string, peerPubKeyHash string, pubKeyHash string, encodedPeerPubKey string) error {
  payloadJoined := strings.Join([]string{pubKeyHash, peerPubKeyHash, encodedPeerPubKey}, ";")

  signatureBytes, err := base64.StdEncoding.DecodeString(signatureEncoded)
  // verify signature
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  pssh := newhash.New()
  pssh.Write([]byte(payloadJoined))
  hashed := pssh.Sum(nil)

  err = rsa.VerifyPSS(peerPubKey, newhash, hashed, signatureBytes, &opts)
  if err != nil {
    mlog.GetLogger().Warn("Signature verification failed")
  }
  return err
}
