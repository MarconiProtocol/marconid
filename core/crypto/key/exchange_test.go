package mcrypto_key

import (
  "crypto/rsa"
  "net/http"
  "sync"
  "testing"
)

func TestKeyManager_registerRPCHandlers(t *testing.T) {
  type fields struct {
    baseKey        *rsa.PrivateKey
    peerPublicKeys *map[string]*rsa.PublicKey
    sentRPC        map[string]bool
    sentRPCMutex   sync.Mutex
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      km := &KeyManager{
        baseKey:        tt.fields.baseKey,
        peerPublicKeys: tt.fields.peerPublicKeys,
        sentRPC:        tt.fields.sentRPC,
        sentRPCMutex:   tt.fields.sentRPCMutex,
      }
      km.registerRPCHandlers()
    })
  }
}

func TestKeyManager_InitiatePublicKeyExchange(t *testing.T) {
  type fields struct {
    baseKey        *rsa.PrivateKey
    peerPublicKeys *map[string]*rsa.PublicKey
    sentRPC        map[string]bool
    sentRPCMutex   sync.Mutex
  }
  type args struct {
    peerIp         string
    peerPubKeyHash string
  }
  tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      km := &KeyManager{
        baseKey:        tt.fields.baseKey,
        peerPublicKeys: tt.fields.peerPublicKeys,
        sentRPC:        tt.fields.sentRPC,
        sentRPCMutex:   tt.fields.sentRPCMutex,
      }
      if err := km.InitiatePublicKeyExchange(tt.args.peerIp, tt.args.peerPubKeyHash); (err != nil) != tt.wantErr {
        t.Errorf("KeyManager.InitiatePublicKeyExchange() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestKeyManager_SendPubKeyExchangeSynRPC(t *testing.T) {
  type fields struct {
    baseKey        *rsa.PrivateKey
    peerPublicKeys *map[string]*rsa.PublicKey
    sentRPC        map[string]bool
    sentRPCMutex   sync.Mutex
  }
  type args struct {
    targetHost     string
    peerPubKeyHash string
  }
  tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      km := &KeyManager{
        baseKey:        tt.fields.baseKey,
        peerPublicKeys: tt.fields.peerPublicKeys,
        sentRPC:        tt.fields.sentRPC,
        sentRPCMutex:   tt.fields.sentRPCMutex,
      }
      if err := km.SendPubKeyExchangeSynRPC(tt.args.targetHost, tt.args.peerPubKeyHash); (err != nil) != tt.wantErr {
        t.Errorf("KeyManager.SendPubKeyExchangeSynRPC() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestKeyManager_handlePubKeyExchangeSynRPC(t *testing.T) {
  type fields struct {
    baseKey        *rsa.PrivateKey
    peerPublicKeys *map[string]*rsa.PublicKey
    sentRPC        map[string]bool
    sentRPCMutex   sync.Mutex
  }
  type args struct {
    r           *http.Request
    w           http.ResponseWriter
    reqInfohash string
    reqPayload  string
  }
  tests := []struct {
    name   string
    fields fields
    args   args
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      km := &KeyManager{
        baseKey:        tt.fields.baseKey,
        peerPublicKeys: tt.fields.peerPublicKeys,
        sentRPC:        tt.fields.sentRPC,
        sentRPCMutex:   tt.fields.sentRPCMutex,
      }
      km.handlePubKeyExchangeSynRPC(tt.args.r, tt.args.w, tt.args.reqPayload)
    })
  }
}

func Test_buildPubKeyExchangePayload(t *testing.T) {
  type args struct {
    peerPubKeyHash string
    pubKey         *rsa.PublicKey
    privateKey     *rsa.PrivateKey
  }
  tests := []struct {
    name string
    args args
    want string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := buildPubKeyExchangePayload(tt.args.peerPubKeyHash, tt.args.pubKey, tt.args.privateKey); got != tt.want {
        t.Errorf("buildPubKeyExchangePayload() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_parsePubKeyExchangePayload(t *testing.T) {
  type args struct {
    payload string
  }
  tests := []struct {
    name               string
    args               args
    wantSignature      string
    wantPubKeyHash     string
    wantPeerPubKeyHash string
    wantEncodedPubKey  string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotSignature, gotPubKeyHash, gotPeerPubKeyHash, gotEncodedPubKey := parsePubKeyExchangePayload(tt.args.payload)
      if gotSignature != tt.wantSignature {
        t.Errorf("parsePubKeyExchangePayload() gotSignature = %v, want %v", gotSignature, tt.wantSignature)
      }
      if gotPubKeyHash != tt.wantPubKeyHash {
        t.Errorf("parsePubKeyExchangePayload() gotPubKeyHash = %v, want %v", gotPubKeyHash, tt.wantPubKeyHash)
      }
      if gotPeerPubKeyHash != tt.wantPeerPubKeyHash {
        t.Errorf("parsePubKeyExchangePayload() gotPeerPubKeyHash = %v, want %v", gotPeerPubKeyHash, tt.wantPeerPubKeyHash)
      }
      if gotEncodedPubKey != tt.wantEncodedPubKey {
        t.Errorf("parsePubKeyExchangePayload() gotEncodedPubKey = %v, want %v", gotEncodedPubKey, tt.wantEncodedPubKey)
      }
    })
  }
}

func TestKeyManager_isRpcSenderReceiverValid(t *testing.T) {
  type fields struct {
    baseKey        *rsa.PrivateKey
    peerPublicKeys *map[string]*rsa.PublicKey
    sentRPC        map[string]bool
    sentRPCMutex   sync.Mutex
  }
  type args struct {
    peerPubKeyHash string
    pubKeyHash     string
  }
  tests := []struct {
    name   string
    fields fields
    args   args
    want   bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      km := &KeyManager{
        baseKey:        tt.fields.baseKey,
        peerPublicKeys: tt.fields.peerPublicKeys,
        sentRPC:        tt.fields.sentRPC,
        sentRPCMutex:   tt.fields.sentRPCMutex,
      }
      if got := km.isRpcSenderReceiverValid(tt.args.peerPubKeyHash, tt.args.pubKeyHash); got != tt.want {
        t.Errorf("KeyManager.isRpcSenderReceiverValid() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_createSignature(t *testing.T) {
  type args struct {
    privateKey     *rsa.PrivateKey
    peerPubKeyHash string
    pubKeyHash     string
    encodedPubKey  string
  }
  tests := []struct {
    name          string
    args          args
    wantSignature string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if gotSignature := createSignature(tt.args.privateKey, tt.args.peerPubKeyHash, tt.args.pubKeyHash, tt.args.encodedPubKey); gotSignature != tt.wantSignature {
        t.Errorf("createSignature() = %v, want %v", gotSignature, tt.wantSignature)
      }
    })
  }
}

func Test_verifySignature(t *testing.T) {
  type args struct {
    peerPubKey        *rsa.PublicKey
    signatureEncoded  string
    peerPubKeyHash    string
    pubKeyHash        string
    encodedPeerPubKey string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := verifySignature(tt.args.peerPubKey, tt.args.signatureEncoded, tt.args.peerPubKeyHash, tt.args.pubKeyHash, tt.args.encodedPeerPubKey); (err != nil) != tt.wantErr {
        t.Errorf("verifySignature() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
