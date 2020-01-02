package mcrypto_key

import (
  "../../../util"
  "../../peer"
  "../../rpc"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "encoding/base64"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "net/http"
  "strings"
)

type PubKeyExchangeSYNArgs struct {
  Payload string
}

type PubKeyExchangeSYNReply struct {
  Payload string
}

type PubKeyExchangeService struct{}

func (dhm *PubKeyExchangeService) PubKeyExchangeSYNRPC(r *http.Request, args *PubKeyExchangeSYNArgs, reply *PubKeyExchangeSYNReply) error {
  return KeyManagerInstance().HandlePubKeyExchangeSYNRPC(r, args, reply)
}

func (km *KeyManager) registerRPCHandlers() {
  pubKeyExchangService := new(PubKeyExchangeService)
  if err := rpc.RegisterService(pubKeyExchangService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for KeyManager")
  }
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

  args := PubKeyExchangeSYNArgs{Payload: payload}
  reply := PubKeyExchangeSYNReply{}
  cli := rpc.NewRPCClient(targetHost)
  if err := cli.Call("PubKeyExchangeService.PubKeyExchangeSYNRPC", &args, &reply); err != nil {
    return errors.New(fmt.Sprintf("SendPubKeyExchangeSynRPC - Nil/error response from %s ; response %v", targetHost, err))
  }

  mlog.GetLogger().Debug(fmt.Sprintf("SendPubKeyExchangeSynRPC - Received response from %s", targetHost))
  signature, pubKeyHash, peerPubKeyHash, encodedPeerPubKey := parsePubKeyExchangePayload(reply.Payload)
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

func (km *KeyManager) HandlePubKeyExchangeSYNRPC(r *http.Request, args *PubKeyExchangeSYNArgs, reply *PubKeyExchangeSYNReply) error {
  // Using the sender IP in the http request; sender could be behind NAT
  ip, _, err := rpc.ParseIPAndPort(r.RemoteAddr)
  if err != nil {
    mlog.GetLogger().Error("PubKeyExchange::handlePubKeyExchangeSynRPC - could not parse remote address / port from request: ", err)
    return err
  }

  signature, pubKeyHash, peerPubKeyHash, encodedPeerPubKey := parsePubKeyExchangePayload(args.Payload)
  if !km.isRpcSenderReceiverValid(peerPubKeyHash, pubKeyHash) {
    mlog.GetLogger().Debug("PubKeyExchange::handlePubKeyExchangeSynRPC - invalid rpc sender, sender is not a peer")
    return errors.New("not recognized as a peer")
  }
  mlog.GetLogger().Debug("PubKeyExchange::handlePubKeyExchangeSynRPC - RECEIVED SYN from", r.RemoteAddr)

  // Decode peerPubKey and verify signature
  peerPubKey := decodePublicKeyB64String(encodedPeerPubKey)
  err = verifySignature(peerPubKey, signature, peerPubKeyHash, pubKeyHash, encodedPeerPubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Signature verification failed: ", err)
  }

  mpeer.PeerManagerInstance().UpdatePeer(peerPubKeyHash, ip)

  // Store public key of peer
  km.SetPeerPublicKey(peerPubKeyHash, peerPubKey)

  mlog.GetLogger().Debug("PubKeyExchange::SendPubKeyExchangeAckRPC - SENDING ACK to", ip)
  reply.Payload = buildPubKeyExchangePayload(peerPubKeyHash, km.GetBasePublicKey(), km.GetBasePrivateKey())
  return nil
}

func buildPubKeyExchangePayload(peerPubKeyHash string, pubKey *rsa.PublicKey, privateKey *rsa.PrivateKey) string {
  // craft our payload
  pubKeyHash, err := mutil.GetInfohashByPubKey(pubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to get pubkeyhash from pub key: ", err)
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
    mlog.GetLogger().Error("Failed to decode pk exchange payload: ", err)
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
    mlog.GetLogger().Fatal("Failed to get pubkeyhash from pub key: ", err)
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
    mlog.GetLogger().Fatal("Failed to generate signature: ", err)
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
    mlog.GetLogger().Warn("Signature verification failed: ", err)
  }
  return err
}
