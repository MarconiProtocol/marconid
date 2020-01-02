package mcrypto_dh

import (
  "errors"
  "fmt"
  "github.com/MarconiProtocol/kyber/group/edwards25519"
  mlog "github.com/MarconiProtocol/log"
  "sync"
)

// TODO: Refactor such that the peer object will be the owner of its own DHKeyInfo. The management logic can be merged into PeerManager

type DHExchangeManager struct {
  suite  *edwards25519.SuiteEd25519
  dhKeys *map[string]*DHKeyInfo

  dhKeysMutex sync.Mutex
}

var instance *DHExchangeManager
var once sync.Once

func DHExchangeManagerInstance() *DHExchangeManager {
  once.Do(func() {
    instance = &DHExchangeManager{}
    instance.initialize()
  })
  return instance
}

func (dhm *DHExchangeManager) initialize() {
  dhm.dhKeys = &map[string]*DHKeyInfo{}
  dhm.suite = edwards25519.NewBlakeSHA256Ed25519()
  dhm.registerRPCHandlers()
}

func (dhm *DHExchangeManager) GetDHKeyInfo(peerPubKeyHash string) (*DHKeyInfo, error) {
  if dhKeyInfo, exists := (*dhm.dhKeys)[peerPubKeyHash]; exists {
    return dhKeyInfo, nil
  }
  return nil, errors.New(fmt.Sprintf("DHExchangeManager - could not get DHKeyInfo for pubkeyhash %s", peerPubKeyHash))
}

func (dhm *DHExchangeManager) RemoveDHKeyInfo(peerPubKeyHash string) {
  mlog.GetLogger().Info("RemoveDHKeyInfo", peerPubKeyHash)
  dhm.dhKeysMutex.Lock()
  if _, exists := (*dhm.dhKeys)[peerPubKeyHash]; exists {
    delete(*dhm.dhKeys, peerPubKeyHash)
  }
  dhm.dhKeysMutex.Unlock()
}
