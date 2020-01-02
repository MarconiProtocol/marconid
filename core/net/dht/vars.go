package mnet_dht

import (
  dht_lib "github.com/MarconiProtocol/dht"
  mlog "github.com/MarconiProtocol/log"
  "sync"
)

const PEER_REQUEST_RESPONSE = "peer_response"
const EDGE_REQUEST_RESPONSE = "edge_response"

type PeerChan chan<- string

type ChanSet struct {
  set map[PeerChan]bool
}

func (set *ChanSet) Add(ch PeerChan) bool {
  _, found := set.set[ch]
  set.set[ch] = true
  return !found
}
func (set *ChanSet) Remove(ch PeerChan) {
  delete(set.set, ch)
}

type MDHT struct {
  DHT         *dht_lib.DHT
  subscribers map[dht_lib.InfoHash]*ChanSet
  subMutex    *sync.Mutex

  //dynamic DHT router ip
  DiscoveredDHTRouterIP *map[string]int
  DiscoveredInfohash    map[string]int
  AllowedPubKey         *map[string]int

  PeerRouteActionHandlerMapping map[string]func(map[string]string)

  Log *mlog.Mlog
  sync.Mutex
}

type Config struct {
  DHTPort   int
  SeedNodes string
}
