package mnet_dht

import (
  "../../../util"
  "../../config"
  "crypto/rsa"
  "github.com/MarconiProtocol/log"
  "sync"
)

type BeaconManager struct {
  baseBeaconMutex  sync.Mutex
  peerBeaconMutex  sync.Mutex
  edgeBeaconMutex  sync.Mutex
  peerRequestMutex sync.Mutex

  baseBeaconDHT    *MDHT
  baseBeaconSignal *chan bool

  edgeBeaconDHT      *MDHT
  edgeBeaconSignal   *chan bool
  edgeRequestStarted bool
  edgeRequestSignal  *chan bool

  peerBeaconDHT      *MDHT
  peerBeaconSignal   *chan bool
  peerRequestStarted *map[string]bool
  peerRequestSignals *map[string]*chan bool // pubkeyhash -> kill channel for peer request goroutines
}

const BEACON_INTERVAL = 5

var beaconManager *BeaconManager
var once sync.Once

// returns the singleton to BeaconManager
func GetBeaconManager() *BeaconManager {
  once.Do(func() {
    beaconManager = &BeaconManager{}

    baseSignal := make(chan bool)
    beaconManager.baseBeaconSignal = &baseSignal

    edgeSignal := make(chan bool)
    beaconManager.edgeBeaconSignal = &edgeSignal
    edgeRequestSignal := make(chan bool)
    beaconManager.edgeRequestSignal = &edgeRequestSignal

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

  mlog.GetLogger().Debug("BeaconManager::StartBaseRouteBeacon - Starting BaseRouteBeacon")
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
  go bm.baseBeaconDHT.AnnouncementBeacon(bm.baseBeaconSignal, baseRouteBeaconKey, 120, mconfig.GetAppConfig().DHT.Announce_Base_Interval_Seconds)
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

// ===== Edge Route Beacon =====
/*
  Creates a new beacon for edge route announcement and requests.
*/
func (bm *BeaconManager) CreateEdgeRouteBeacon(actionFunc func(map[string]string)) {
  if bm.edgeBeaconDHT != nil {
    mlog.GetLogger().Debug("BeaconManager::StartEdgeRouteBeacon - Tried to start another edge route beacon while one is active, no-op")
    return
  }

  mlog.GetLogger().Debug("BeaconManager::StartEdgeRouteBeacon - Starting EdgeRouteBeacon")
  conf := &Config{
    DHTPort:   24803,
    SeedNodes: mconfig.GetAppConfig().DHT.BootNodes,
  }

  var err error
  if bm.edgeBeaconDHT, err = NewMDHT(conf); err != nil {
    mlog.GetLogger().Fatal("BeaconManager::StartEdgeRouteBeacon - Error: Failed to start dht instance for edge route beacon", err)
  }

  bm.edgeBeaconDHT.PeerRouteActionHandlerMapping[EDGE_REQUEST_RESPONSE] = actionFunc
  go bm.edgeBeaconDHT.handleEdgeRequestResponses()
}

/*
  Announce on the edge route beacon that this node is a part of the edge route
*/
func (bm *BeaconManager) StartEdgeRouteAnnouncement(edgeRouteBeaconKey *rsa.PublicKey) {
  if bm.edgeBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StartEdgeRouteAnnouncement - Tried to start edge route announcement when no beacon was created")
    return
  }
  go bm.edgeBeaconDHT.AnnouncementBeacon(bm.edgeBeaconSignal, edgeRouteBeaconKey, 0, mconfig.GetAppConfig().DHT.Announce_Base_Interval_Seconds)
}

/*
  Start requesting for responses from peers on the base route
*/
func (bm *BeaconManager) StartEdgeRouteRequest(baseRouteBeaconKey *rsa.PublicKey) {
  if bm.isEdgeRouteRequestStarted() {
    mlog.GetLogger().Debug("StartBaseRouteRequest is already started, no-op")
    return
  }

  baseRoutePubKeyHash, err := mutil.GetInfohashByPubKey(baseRouteBeaconKey)
  if err != nil {
    mlog.GetLogger().Fatal("BeaconManager::StartBaseRouteRequest - Failed to get pkh of base route key")
  }
  bm.edgeRequestStarted = true

  go bm.edgeBeaconDHT.RequestBeacon(bm.edgeRequestSignal, baseRoutePubKeyHash, "", []string{}, BEACON_INTERVAL)
}

/*
  Stop requesting for responses from peers on the base route
*/
func (bm *BeaconManager) StopEdgeRouteRequest() {
  if !bm.isEdgeRouteRequestStarted() {
    mlog.GetLogger().Debug("BaseRouteRequest is not started, no-op")
    return
  }

  bm.edgeRequestStarted = false
  *bm.edgeRequestSignal <- true
}

/*
  Check if the edge route request has been started
*/
func (bm *BeaconManager) isEdgeRouteRequestStarted() bool {
  return bm.edgeRequestStarted
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
  bm.peerBeaconDHT.PeerRouteActionHandlerMapping[PEER_REQUEST_RESPONSE] = actionFunc
  go bm.peerBeaconDHT.handlePeerRequestResponses()
}

/*
  Start announcing this node's PubKeyHash on the peer route
*/
func (bm *BeaconManager) StartPeerRouteAnnouncement() {
  if bm.baseBeaconDHT == nil {
    mlog.GetLogger().Debug("BeaconManager::StartPeerRouteAnnouncement - Tried to start peer route announcement when no beacon was created")
    return
  }

  mlog.GetLogger().Debug("BeaconManager::StartPeerRouteAnnouncement")
  go bm.peerBeaconDHT.SelfAnnouncementBeacon(bm.peerBeaconSignal, 0, mconfig.GetAppConfig().DHT.Announce_Self_Interval_Seconds)
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
  go bm.peerBeaconDHT.RequestBeacon(&signal, peerPubKeyHash, "", []string{}, mconfig.GetAppConfig().DHT.Request_Peers_Interval_Seconds)
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
