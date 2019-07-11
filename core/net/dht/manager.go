package mnet_dht

import (
  "../../config"
  "crypto/rsa"
  "gitlab.neji.vm.tc/marconi/log"
  "sync"
)

type BeaconManager struct {
  baseBeaconMutex  sync.Mutex
  peerBeaconMutex  sync.Mutex
  peerRequestMutex sync.Mutex

  baseBeaconDHT    *MDHT
  baseBeaconSignal *chan bool

  peerBeaconDHT    *MDHT
  peerBeaconSignal *chan bool

  peerRequestStarted *map[string]bool
  peerRequestSignals *map[string]*chan bool // pubkeyhash -> kill channel for peer request goroutines
}

var beaconManager *BeaconManager
var once sync.Once

// returns the singleton to BeaconManager
func GetBeaconManager() *BeaconManager {
  once.Do(func() {
    beaconManager = &BeaconManager{}

    baseSignal := make(chan bool)
    beaconManager.baseBeaconSignal = &baseSignal

    peerSignal := make(chan bool)
    beaconManager.peerBeaconSignal = &peerSignal

    peerRequestSignals := make(map[string]*chan bool)
    beaconManager.peerRequestSignals = &peerRequestSignals

    peerRequestStarted := make(map[string]bool)
    beaconManager.peerRequestStarted = &peerRequestStarted
  })
  return beaconManager
}

// ===== Base Route Beacon =====
/*
  Creates a new beacon for base route announcement and requests.
*/
func (bm *BeaconManager) CreateBaseRouteBeacon() {
  if bm.baseBeaconDHT != nil {
    mlog.GetLogger().Debug("BeaconManager::StartBaseRouteBeacon - Tried to start another base route beacon while one is active, no-op")
    return
  }

  mlog.GetLogger().Debug("BeaconManager::StartBaseRouteBeacon - Starting PeerRouteBeacon")
  conf := &Config{
    DHTPort:   24801,
    SeedNodes: mconfig.GetAppConfig().DHT.BootNodes,
  }

  var err error
  if bm.baseBeaconDHT, err = NewMDHT(conf); err != nil {
    mlog.GetLogger().Fatal("BeaconManager::StartBaseRouteBeacon - Error: Failed to start dht instance for base route beacon", err)
  }
}

/*
  Announce on the base route beacon that this node is a part of the base route
*/
func (bm *BeaconManager) StartBaseRouteAnnouncement(baseRouteBeaconKey *rsa.PublicKey) {
  if bm.baseBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StartBaseRouteAnnouncement - Tried to start base route announcement when no beacon was created")
    return
  }

  // TODO probably need to lock the read/write for this
  go bm.baseBeaconDHT.AnnouncementBeacon(bm.baseBeaconSignal, baseRouteBeaconKey, 120, mconfig.GetAppConfig().DHT.AnnounceSelfIntervalSeconds)
}

/*
  Stop announcing on the base route
*/
func (bm *BeaconManager) StopBaseRouteAnnouncement() {
  if bm.baseBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StopBaseRouteBeacon - Tried to stop base route announcement while none is active, no-op")
    return
  }
  *bm.baseBeaconSignal <- true
}

// ===== Peer Route Beacon =====
/*
  Creates a new beacon for peer route announcement and requests.
*/
func (bm *BeaconManager) CreatePeerRouteBeacon(actionFunc func(map[string]string)) {
  if bm.peerBeaconDHT != nil {
    mlog.GetLogger().Debug("BeaconManager::StartPeerRouteBeacon - Tried to start another peer route beacon while one is active, no-op")
    return
  }

  conf := &Config{
    DHTPort:   24800,
    SeedNodes: mconfig.GetAppConfig().DHT.BootNodes,
  }

  var err error
  if bm.peerBeaconDHT, err = NewMDHT(conf); err != nil {
    mlog.GetLogger().Fatal("BeaconManager - Error: Failed to start dht instance for peer route beacon", err)
  }
  // Set up the callback for when a PeerRouteRequest has a response
  bm.peerBeaconDHT.PeerRouteActionList[PEER_REQUEST_RESPONSE] = actionFunc
  go bm.peerBeaconDHT.handlePeerRequestResponses()
}

/*
  Start announcing this node's PKH on the peer route
*/
func (bm *BeaconManager) StartPeerRouteAnnouncement() {
  if bm.baseBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StartPeerRouteAnnouncement - Tried to start peer route announcement when no beacon was created")
    return
  }

  mlog.GetLogger().Debug("BeaconManager::StartPeerRouteAnnouncement - peerRouteKey: ")
  go bm.peerBeaconDHT.SelfAnnouncementBeacon(bm.peerBeaconSignal, 0, 2)
}

/*
  Stop announcing on the peer route
*/
func (bm *BeaconManager) StopPeerRouteAnnouncement() {
  if bm.peerBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StopPeerRouteAnnouncement - Tried to stop peer route announcement while none is active, no-op")
    return
  }
  *bm.peerBeaconSignal <- true
}

/*
Start requesting for peerPubKeyHash on the peer route
*/
func (bm *BeaconManager) StartPeerRouteRequest(peerPubKeyHash string) {
  if bm.isPeerRouteRequestStarted(peerPubKeyHash) {
    mlog.GetLogger().Debug("PeerRouteRequest is already started, no-op")
    return
  }

  signal := make(chan bool)
  bm.peerRequestMutex.Lock()
  (*bm.peerRequestStarted)[peerPubKeyHash] = true
  (*bm.peerRequestSignals)[peerPubKeyHash] = &signal
  go bm.peerBeaconDHT.RequestBeacon(&signal, peerPubKeyHash, "", []string{}, 360000, 5)
  bm.peerRequestMutex.Unlock()
}

func (bm *BeaconManager) StopPeerRouteRequest(peerPubKeyHash string) {
  if !bm.isPeerRouteRequestStarted(peerPubKeyHash) {
    mlog.GetLogger().Debug("PeerRouteRequest is not started, no-op")
    return
  }

  bm.peerRequestMutex.Lock()
  signal := (*bm.peerRequestSignals)[peerPubKeyHash]
  (*bm.peerRequestStarted)[peerPubKeyHash] = false
  *signal <- true
  bm.peerRequestMutex.Unlock()
}

func (bm *BeaconManager) isPeerRouteRequestStarted(peerPubKeyHash string) bool {
  if started, exists := (*bm.peerRequestStarted)[peerPubKeyHash]; exists {
    if started {
      return true
    }
  }
  return false
}
