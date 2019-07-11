package mcrypto_dh

import (
  "../key"
  "../../peer"
  "../../rpc"
  "../../../util"
  "bytes"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "encoding/base64"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
  "net/http"
  "strings"
)

func (dhm *DHExchangeManager) registerRPCHandlers() {
  // Register functions to handle named RPCs
  rpc.RegisterRpcHandler(rpc.REQUEST_DH_KEY_EXCHANGE_SYN, dhm.handleSynRPC)
}

func (dhm *DHExchangeManager) InitiateDHKeyExchange(peerPubKeyHash string) error {
  // Get the peer's ip from peer manager
  peer, err := mpeer.PeerManagerInstance().GetPeer(peerPubKeyHash)
  if err != nil {
    return errors.New(fmt.Sprintf("InitiateDHKeyExchange - No peer found for %v", peerPubKeyHash))
  }

  host := peer.Ip
  if host == "" {
    return errors.New(fmt.Sprintf("InitiateDHKeyExchange - Peer IP was not set %v", peerPubKeyHash))
  }

  peerPubKey, err := mcrypto_key.KeyManagerInstance().GetPeerPublicKey(peerPubKeyHash)
  if err != nil {
    return errors.New(fmt.Sprintf("InitiateDHKeyExchange: No key found for peer, error case, shouldn't call InitiateDHKeyExchange before having pub key %v", peerPubKeyHash))
  }

  dhKeyInfo := dhm.getDHKeyInfoForPeer(peerPubKeyHash)
  if dhKeyInfo.SymmetricKey == nil {
    return dhm.SendSynRPC(host, peerPubKey, dhKeyInfo)
  }

  return nil
}

/*
  SendSynRPC will send out the RPC to start the DH exchange handshake with the specific peer
*/
func (dhm *DHExchangeManager) SendSynRPC(targetHost string, peerPubKey *rsa.PublicKey, dhKeyInfo *DHKeyInfo) error {
  mlog.GetLogger().Debug(fmt.Sprintf("DHExchange::SendSynRPC - Sending syn to %s", targetHost))
  payload := buildPayload(peerPubKey, dhKeyInfo)

  response := rpc.SendRPC(targetHost, rpc.RPC_PORT, rpc.REQUEST_DH_KEY_EXCHANGE_SYN, payload)
  if response.Error != "" {
    return errors.New(fmt.Sprintf("DHExchange::SendSynRPC - Nil/error response from: %s; response: %v", targetHost, response.Error))
  }

  mlog.GetLogger().Debug(fmt.Sprintf("DHExchange::SendSynRPC - Received response from: %s", targetHost))

  senderPubKeyHash, encodedPeerDHPubKey := parsePayload(response.Result)

  // Try to decode the payload into the peerpublickey
  senderDHhKeyInfo := dhm.getDHKeyInfoForPeer(senderPubKeyHash)
  err := dhKeyInfo.DecodeStringToPeerPublicKey(encodedPeerDHPubKey)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to decode peer public key %s", err))
  }

  senderDHhKeyInfo.GenDHSymmetricKey()
  return nil
}

/*
  handleSynRPC is invoked in response to receiving rpc.REQUEST_DH_KEY_EXCHANGE_SYN
  The symmetric key for the peer is generated and an ACK is sent in response
*/
func (dhm *DHExchangeManager) handleSynRPC(r *http.Request, w http.ResponseWriter, reqInfohash string, reqPayload string) {
  // Using the sender IP in the http request; sender could be behind NAT
  ip, _, err := rpc.ParseIPAndPort(r.RemoteAddr)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("DHExchange::handleSynRPC - could not parse remote address / port from request"))
    err = rpc.WriteResponseWithError(&w, "could not parse remote address from request")
    return
  }

  senderPubKeyHash, encodedPeerDHPubKey := parsePayload(reqPayload)
  mlog.GetLogger().Debug(fmt.Sprintf("DHExchange::handleSynRPC : RECEIVED SYN from %s", r.RemoteAddr))
  senderPubKey, err := mcrypto_key.KeyManagerInstance().GetPeerPublicKey(senderPubKeyHash)
  if err != nil {
    mlog.GetLogger().Debug("Don't have the public key of the peer that sent DHExchange syn, possible error case")
    err = rpc.WriteResponseWithError(&w, "did not have peer public key")
    return
  }

  // Get DHKeyInfo for peer and try to decode peer's public key into the struct instance
  dhKeyInfo := dhm.getDHKeyInfoForPeer(senderPubKeyHash)
  err = dhKeyInfo.DecodeStringToPeerPublicKey(encodedPeerDHPubKey)
  if err == nil {
    // We have our private/public key pair and also received our peer's public key, generate symmetric key
    dhKeyInfo.GenDHSymmetricKey()

    mlog.GetLogger().Debug(fmt.Sprintf("DHExchange::replying with ACK - SENDING ACK to %s", ip))
    payload := buildPayload(senderPubKey, dhKeyInfo)
    err := rpc.WriteResponseWithResult(&w, payload)
    if err != nil {
      mlog.GetLogger().Debug(fmt.Sprintf("DHExchange:: error writing to response, err: %s", err))
    }
  }
}

/*
  Attempts to get the peer's DHKeyInfo, if one does not exist, then it will be generated
*/
func (dhm *DHExchangeManager) getDHKeyInfoForPeer(peerPubKeyHash string) *DHKeyInfo {
  dhm.dhKeysMutex.Lock()
  dhKeyInfo, exists := (*dhm.dhKeys)[peerPubKeyHash]
  if !exists {
    mlog.GetLogger().Debug("| DHExchange: generating dhKeyInfo struct for peer:", peerPubKeyHash)
    dhKeyInfo = dhm.genPrivatePublicKeysForPeer(peerPubKeyHash)
  }
  dhm.dhKeysMutex.Unlock()
  return dhKeyInfo
}

/* Source:
https://en.wikipedia.org/wiki/Elliptic-curve_Diffie%E2%80%93Hellman

The following example will illustrate how a key establishment is made.
Suppose Alice wants to establish a shared key with Bob, but the only channel available for them may be eavesdropped by a third party.

Initially, the domain parameters (that is, (p,a,b,G,n,h) in the prime case or (m,f(x),a,b,G,n,h) in the binary case) must be agreed upon.
In our case this agreement happens in the form of using a specific cryptographic suite, SuiteEd25519
[source on domain parameters: https://en.wikipedia.org/wiki/Elliptic-curve_cryptography#Domain_parameters]

Also, each party must have a key pair suitable for elliptic curve cryptography,
consisting of a private key d (a randomly selected integer in the interval [1,n-1] and a public key represented by a point Q (where Q=dG, that is, the result of adding G to itself d times).
In our usage of the marconi/kyber library a Scalar object is chosen as d, and a Point object chosen as Q using d

Let Alice's key pair be (d_{A},Q_{A}) and Bob's key pair be (d_{B},Q_{B}). Each party must know the other party's public key prior to execution of the protocol.
Alice computes point (x_{k},y_{k})=d_{A}Q_{B}. Bob computes point (x_{k},y_{k})=d_{B}Q_{A}. The shared secret is x_{k} (the x coordinate of the point).

Most standardized protocols based on ECDH derive a symmetric key from x_{k} using some hash-based key derivation function.

The shared secret calculated by both parties is equal, because d_{A}Q_{B}=d_{A}d_{B}G=d_{B}d_{A}G=d_{B}Q_{A}.

The only information about her private key that Alice initially exposes is her public key.
So, no party other than Alice can determine Alice's private key, unless that party can solve the elliptic curve discrete logarithm problem.
Bob's private key is similarly secure. No party other than Alice or Bob can compute the shared secret, unless that party can solve the elliptic curve Diffie–Hellman problem.
*/
func (dhm *DHExchangeManager) genPrivatePublicKeysForPeer(peerPubKeyHash string) *DHKeyInfo {
  // used as a source of cryptographic randomness
  cypherStream := dhm.suite.RandomStream()

  // generate the private key (a scalar used to pick a point on the curve)
  private := dhm.suite.Scalar().Pick(cypherStream)
  // pick the point on the curve based on the generated scalar
  public := dhm.suite.Point().Mul(private, nil)

  dhKeyInfo := &DHKeyInfo{}
  dhKeyInfo.PrivateKey = private
  dhKeyInfo.PublicKey = public
  dhKeyInfo.SymmetricKeyBytes = new(bytes.Buffer)
  keySignal := make(chan *bytes.Buffer)
  dhKeyInfo.KeySignal = &keySignal

  (*dhm.dhKeys)[peerPubKeyHash] = dhKeyInfo
  return dhKeyInfo
}

/*
  buildPayload is used to generate the rpc payload from the provided arguments
*/
func buildPayload(peerPubKey *rsa.PublicKey, dhKeyInfo *DHKeyInfo) string {
  encodedPublicKey := dhKeyInfo.EncodePublicKeyToString()

  // data payload format
  // <encodedPublicKey>
  dataPayload := encodedPublicKey

  pubKeyHash, err := mutil.GetInfohashByPubKey(mcrypto_key.KeyManagerInstance().GetBasePublicKey())
  if err != nil {
    mlog.GetLogger().Fatalf("KeyManager:: Error loading Key at [%v], ", err)
  }

  // payload format
  // <peerPubKeyHash>;<signature>;<cipher>
  signature, cipherText := encryptDataPayload(peerPubKey, dataPayload)
  payload := strings.Join([]string{pubKeyHash, signature, cipherText}, ";")

  // base64 encode the payload
  payload = base64.StdEncoding.EncodeToString([]byte(payload))

  return payload
}

/*
  parsePayload accepts the payload from the rpc call, and spits out the info parsed ( and the data decrypted )
*/
func parsePayload(payload string) (peerPubKeyHash string, encodedPeerDHKey string) {
  // base64 decode the payload
  payloadBytes, err := base64.StdEncoding.DecodeString(payload)
  if err != nil {
    mlog.GetLogger().Error("Failed to decode dh exchange payload")
    return "", ""
  }
  payload = string(payloadBytes)

  // split the string payload into the peerPubKeyHash, signature and encryptedPayload
  payloadSplit := strings.Split(payload, ";")
  peerPubKeyHash, signature, encryptedDataPayload := payloadSplit[0], payloadSplit[1], payloadSplit[2]

  // decrypt the encryptedDataPayload
  peerPubKey, err := mcrypto_key.KeyManagerInstance().GetPeerPublicKey(peerPubKeyHash)
  if err != nil {
    mlog.GetLogger().Infof("DHExchange::parsePayload: No key found for peer %v", peerPubKeyHash)
    // TODO: ayuen separate out the functions
    // its possible we dont have the key at this point
    return "", ""
  }

  plainTextDataPayload := decryptDataPayload(peerPubKey, signature, encryptedDataPayload)
  encodedPeerDHKey = plainTextDataPayload

  return peerPubKeyHash, encodedPeerDHKey
}

/*
  encryptDataPayload encrypts the provided data with the provided public key, and generates a signature and cipher text
*/
func encryptDataPayload(peerPubKey *rsa.PublicKey, dataPayload string) (signature string, cipherText string) {
  // encrypt the data payload with peer's public key
  hash := sha256.New()
  dataPayloadBytes := []byte(dataPayload)
  labelBytes := []byte("")

  cipherTextBytes, err := rsa.EncryptOAEP(hash, rand.Reader, peerPubKey, dataPayloadBytes, labelBytes)
  if err != nil {
    mlog.GetLogger().Fatal("Encryption of data payload failed")
  }

  // generate signature
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  PSSmessage := dataPayloadBytes
  newhash := crypto.SHA256
  pssh := newhash.New()
  pssh.Write(PSSmessage)
  hashed := pssh.Sum(nil)

  privateKey := mcrypto_key.KeyManagerInstance().GetBasePrivateKey()

  signatureBytes, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, &opts)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to generate signature")
  }

  return base64.StdEncoding.EncodeToString(signatureBytes), base64.StdEncoding.EncodeToString(cipherTextBytes)
}

/*
  decryptDataPayload decrypts the provided data with the provided key and verifies integrity using the provided signature
*/
func decryptDataPayload(peerPubKey *rsa.PublicKey, signatureEncoded string, cipherTextEncoded string) string {

  // decode the payload and signature
  signatureBytes, err := base64.StdEncoding.DecodeString(signatureEncoded)
  cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherTextEncoded)

  // decrypt the payload
  hash := sha256.New()
  labelBytes := []byte("")

  privateKey := mcrypto_key.KeyManagerInstance().GetBasePrivateKey()

  // function comment block says the random is used to prevent side channel timing attacks, thats awesome
  plainTextBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, cipherTextBytes, labelBytes)
  if err != nil {
    mlog.GetLogger().Fatal("Decryption of data payload failed")
  }

  // verify signature
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  PSSmessage := plainTextBytes
  newhash := crypto.SHA256
  pssh := newhash.New()

  pssh.Write(PSSmessage)
  hashed := pssh.Sum(nil)

  err = rsa.VerifyPSS(peerPubKey, newhash, hashed, signatureBytes, &opts)
  if err != nil {
    mlog.GetLogger().Fatal("Signature verification failed")
  }

  return string(plainTextBytes)
}
