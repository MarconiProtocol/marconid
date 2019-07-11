package mcrypto_key

import (
  "bytes"
  "crypto/rsa"
  "crypto/x509"
  "encoding/base64"
  "encoding/binary"
  "encoding/pem"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
  "io/ioutil"
  "math/big"
  "strings"
  "sync"

  "../../../util"
)

// TODO: Refactor such that the Peer object will be the owner of its own private/public keys. The management logic can be merged into PeerManager

type KeyManager struct {
  baseKey *rsa.PrivateKey

  peerPublicKeys *map[string]*rsa.PublicKey

  sentRPC      map[string]bool
  sentRPCMutex sync.Mutex
}

var instance *KeyManager
var once sync.Once

const (
  PRIVATE_KEY_PATH = "/opt/marconi/etc/marconid/keys/mpkey"
  PUBLIC_KEY_PATH  = PRIVATE_KEY_PATH + ".pub"
)

func KeyManagerInstance() *KeyManager {
  once.Do(func() {
    instance = &KeyManager{}
    instance.initialize()
  })
  return instance
}

func (km *KeyManager) initialize() {
  km.peerPublicKeys = &map[string]*rsa.PublicKey{}
  km.sentRPC = map[string]bool{}
  km.registerRPCHandlers()
}

func (km *KeyManager) EnsurePrivatePublicKeysGenerated() {
  _, err := readKeyFile(PRIVATE_KEY_PATH)
  if err != nil {
    // no private key file found, assume none exists, generate new private and public keys
    km.generatePrivatePublicKeyPair()
  }
  km.LoadDefaultBaseKey()
}

func (km *KeyManager) LoadDefaultBaseKey() {
  km.LoadBaseKey(PRIVATE_KEY_PATH)
}

func (km *KeyManager) LoadBaseKey(path string) {
  km.baseKey = mutil.LoadKey(path)
}

func (km *KeyManager) GetBasePrivateKey() *rsa.PrivateKey {
  return km.baseKey
}

func (km *KeyManager) GetBasePublicKey() *rsa.PublicKey {
  if km.baseKey == nil {
    return nil
  }
  return &km.baseKey.PublicKey
}

func (km *KeyManager) GetBasePublicKeyHash() string {
  hash, err := mutil.GetInfohashByPubKey(km.GetBasePublicKey())
  if err != nil {
    mlog.GetLogger().Fatal("Could not generate pubkeyhash ", err)
  }
  return hash
}

func (km *KeyManager) SetPeerPublicKey(peerPubKeyHash string, peerPubKey *rsa.PublicKey) {
  // Store to file later as well
  // For now lets just overwrite
  (*km.peerPublicKeys)[peerPubKeyHash] = peerPubKey

  filename := getPeerPubKeyFileName(peerPubKeyHash)
  writeKeyToFile(peerPubKey, filename)
}

func (km *KeyManager) DeletePeerPublicKey(peerPubKeyHash string) {
  if _, exists := (*km.peerPublicKeys)[peerPubKeyHash]; exists {
    delete(*km.peerPublicKeys, peerPubKeyHash)
  }
}

func (km *KeyManager) GetPeerPublicKey(peerPubKeyHash string) (*rsa.PublicKey, error) {
  // First Check if key has been loaded in memory
  if peerPubKey, exists := (*km.peerPublicKeys)[peerPubKeyHash]; exists {
    return peerPubKey, nil
  }
  // Next Check if the key is store on disk
  filename := getPeerPubKeyFileName(peerPubKeyHash)
  if peerPubKey, err := LoadRSAPubKey(filename); err == nil {
    mlog.GetLogger().Infof("Loaded peer pub key from file: %v", peerPubKeyHash)
    return peerPubKey, nil
  }
  return nil, errors.New(fmt.Sprintf("KeyManager::Peer PublicKey does not exist for peer: %v", peerPubKeyHash))
}

func getPeerPubKeyFileName(peerPubKeyHash string) string {
  return peerPubKeyHash + "_key.pub"
}

func encodePubKeyToB64String(pubKey *rsa.PublicKey) string {
  pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Error marshalling SELF public key to pkix bytes")
  }
  return base64.StdEncoding.EncodeToString(pubKeyBytes)
}

func decodePublicKeyB64String(encodedPeerPubKey string) *rsa.PublicKey {
  // Decode peerPubKey string to bytes
  peerPubKeyBytes, err := base64.StdEncoding.DecodeString(encodedPeerPubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Error decoding encodedPubKey")
  }

  // Parse bytes to a public key
  peerPubKeyInterface, err := x509.ParsePKIXPublicKey(peerPubKeyBytes)
  if err != nil {
    mlog.GetLogger().Fatal(err)
  }
  switch peerPubKey := peerPubKeyInterface.(type) {
  case *rsa.PublicKey:
    return peerPubKey
  default:
    mlog.GetLogger().Fatal("Unknown public key type")
  }
  return nil
}

func writeKeyToFile(pubKey *rsa.PublicKey, filename string) {
  pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
  if err != nil {
    mlog.GetLogger().Fatal("Error marshalling public key to pkix bytes")
  }
  pubKeyFileBytes := pem.EncodeToMemory(&pem.Block{
    Type:  "RSA PUBLIC KEY",
    Bytes: pubKeyBytes,
  })
  ioutil.WriteFile(filename, pubKeyFileBytes, 0644)
}

func getRsaValues(data []byte) (format string, e *big.Int, n *big.Int, err error) {
  //fmt.Println("getRsaValues:", data)
  data, length, err := readLength(data)
  if err != nil {
    return
  }

  format = string(data[0:length])
  data = data[length:]

  data, length, err = readLength(data)
  if err != nil {
    return
  }

  data, e, err = readBigInt(data, length)
  if err != nil {
    return
  }

  data, length, err = readLength(data)
  if err != nil {
    return
  }

  data, n, err = readBigInt(data, length)
  if err != nil {
    return
  }
  return
}

func readLength(data []byte) ([]byte, uint32, error) {
  l_buf := data[0:4]
  buf := bytes.NewBuffer(l_buf)
  var length uint32
  err := binary.Read(buf, binary.BigEndian, &length)
  if err != nil {
    return nil, 0, err
  }
  return data[4:], length, nil
}

func readBigInt(data []byte, length uint32) ([]byte, *big.Int, error) {
  var bigint = new(big.Int)
  bigint.SetBytes(data[0:length])
  return data[length:], bigint, nil
}

func LoadRSAPubKey(pubKeyFilePath string) (pubKey *rsa.PublicKey, err error) {
  keyByte, err := ioutil.ReadFile(pubKeyFilePath)
  if err != nil {
    return nil, err
  }
  tokens := strings.Split(string(keyByte), " ")
  if len(tokens) < 2 {
    fmt.Errorf("Invalid key format; must contain at least two fields (keytype data [comment])")
    return nil, err
  }
  //fmt.Println("pubKeyStr/tokens:", tokens)
  key_type := tokens[0]
  data, err := base64.StdEncoding.DecodeString(tokens[1])
  if err != nil {
    fmt.Errorf("failed to decode string %s", tokens[1])
    return nil, err
  }
  format, e, n, err := getRsaValues(data)
  if format != key_type {
    fmt.Errorf("Key type said %s, but encoded format said %s.  These should match!", key_type, format)
    return nil, err
  }
  pubKey = &rsa.PublicKey{
    N: n,
    E: int(e.Int64()),
  }
  return
}
