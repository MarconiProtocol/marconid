package mnet_dht

import (
  "crypto/rsa"
  "reflect"
  "sync"
  "testing"
)

func TestGetBeaconManager(t *testing.T) {
  tests := []struct {
    name string
    want *BeaconManager
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetBeaconManager(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetBeaconManager() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBeaconManager_CreateBaseRouteBeacon(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.CreateBaseRouteBeacon()
    })
  }
}

func TestBeaconManager_StartBaseRouteAnnouncement(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  type args struct {
    baseRouteBeaconKey *rsa.PublicKey
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
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StartBaseRouteAnnouncement(tt.args.baseRouteBeaconKey)
    })
  }
}

func TestBeaconManager_StopBaseRouteAnnouncement(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StopBaseRouteAnnouncement()
    })
  }
}

func TestBeaconManager_CreatePeerRouteBeacon(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  type args struct {
    actionFunc func(map[string]string)
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
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.CreatePeerRouteBeacon(tt.args.actionFunc)
    })
  }
}

func TestBeaconManager_StartPeerRouteAnnouncement(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StartPeerRouteAnnouncement()
    })
  }
}

func TestBeaconManager_StopPeerRouteAnnouncement(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StopPeerRouteAnnouncement()
    })
  }
}

func TestBeaconManager_StartPeerRouteRequest(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  type args struct {
    peerPubKeyHash string
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
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StartPeerRouteRequest(tt.args.peerPubKeyHash)
    })
  }
}

func TestBeaconManager_StopPeerRouteRequest(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  type args struct {
    peerPubKeyHash string
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
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      bm.StopPeerRouteRequest(tt.args.peerPubKeyHash)
    })
  }
}

func TestBeaconManager_isPeerRouteRequestStarted(t *testing.T) {
  type fields struct {
    baseBeaconMutex    sync.Mutex
    peerBeaconMutex    sync.Mutex
    peerRequestMutex   sync.Mutex
    baseBeaconDHT      *MDHT
    baseBeaconSignal   *chan bool
    peerBeaconDHT      *MDHT
    peerBeaconSignal   *chan bool
    peerRequestStarted *map[string]bool
    peerRequestSignals *map[string]*chan bool
  }
  type args struct {
    peerPubKeyHash string
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
      bm := &BeaconManager{
        baseBeaconMutex:    tt.fields.baseBeaconMutex,
        peerBeaconMutex:    tt.fields.peerBeaconMutex,
        peerRequestMutex:   tt.fields.peerRequestMutex,
        baseBeaconDHT:      tt.fields.baseBeaconDHT,
        baseBeaconSignal:   tt.fields.baseBeaconSignal,
        peerBeaconDHT:      tt.fields.peerBeaconDHT,
        peerBeaconSignal:   tt.fields.peerBeaconSignal,
        peerRequestStarted: tt.fields.peerRequestStarted,
        peerRequestSignals: tt.fields.peerRequestSignals,
      }
      if got := bm.isPeerRouteRequestStarted(tt.args.peerPubKeyHash); got != tt.want {
        t.Errorf("BeaconManager.isPeerRouteRequestStarted() = %v, want %v", got, tt.want)
      }
    })
  }
}
