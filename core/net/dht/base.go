package mnet_dht

import (
  "../../config"
  "../../runtime"
  "../../crypto/key"
  "../../../util"
  "crypto/rsa"
  "encoding/hex"
  "fmt"
  dht_lib "gitlab.neji.vm.tc/marconi/dht"
  "gitlab.neji.vm.tc/marconi/log"
  "strconv"
  "strings"
  "sync"
  "time"
)

/*
	Returns a configured MDHT instance based on the provided config
*/
func NewMDHT(conf *Config) (*MDHT, error) {
  // DHT lib
  dhtLibConf := initalizeDHTLibConfig(conf)
  dhtLib, err := dht_lib.New(dhtLibConf)
  if err != nil {
    return nil, err
  }

  // Log
  logger, err := mlog.GetLogInstance("dht")
  if err != nil {
    mlog.GetLogger().Error("Couldn't initialize dht log:", err)
  }

  // Create new MDHT instance
  mdht := &MDHT{
    DHT:                 dhtLib,
    subscribers:         make(map[dht_lib.InfoHash](*ChanSet)),
    subMutex:            &sync.Mutex{},
    PeerRouteActionList: make(map[string]func(map[string]string)),
    Log:                 logger,
  }

  mlog.GetLogger().Info("Creating new DHT with Namespace: ", dhtLibConf.Namespace, ", Router IPs: ", dhtLibConf.DHTRouters, ", Port: ", dhtLibConf.Port)
  if err = mdht.Run(); err != nil {
    mlog.GetLogger().Fatalf("DHT start error: %v", err)
  }

  return mdht, nil
}

/*
	Initialize a config object for the dht lib
*/
func initalizeDHTLibConfig(conf *Config) *dht_lib.Config {
  dhtLibConf := dht_lib.NewConfig()
  dhtLibConf.RateLimit = mconfig.GetAppConfig().DHT.MaxIncomingPacketsPerSecond
  dhtLibConf.ClientPerMinuteLimit = mconfig.GetAppConfig().DHT.MaxPerClientIncomingPacketsPerMinute
  dhtLibConf.SaveRoutingTable = mconfig.GetAppConfig().DHT.CacheRoutingTableToDisk
  dhtLibConf.Address = mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr()
  dhtLibConf.Namespace = strconv.Itoa(mconfig.GetUserConfig().Blockchain.ChainID)
  dhtLibConf.Port = conf.DHTPort
  dhtLibConf.DHTRouters = conf.SeedNodes
  return dhtLibConf
}

/*
	Start the DHT instance
*/
func (m *MDHT) Run() (err error) {
  m.Lock()
  err = m.DHT.Start()
  if err != nil {
    return err
  }
  m.Unlock()
  return
}

// Runs as a goroutine, to handle responses from the dht client
func (m *MDHT) handlePeerRequestResponses() {
  // wait on an item from PeersRequestResults
  for {
    select {
    case responses := <-m.DHT.PeersRequestResults:
      // responses is a map[InfoHash][]string, where the slice contains binary encoded peer address in (ip:port) form
      for infohash, peers := range responses {
        // infohash is the pubkeyhash but encoded as hex, need to convert this back to a string pubkeyhash
        bytes, err := hex.DecodeString(infohash.String())
        if err != nil {
          mlog.GetLogger().Error("Could not decode infohash from dht back to a pubkeyhash")
          continue
        }
        pubkeyhash := string(bytes)
        // Iterate through peers found to have responded for a particular peer pubkeyhash
        for _, peerAddress := range peers {
          ipport := strings.Split(dht_lib.DecodePeerAddress(peerAddress), ":")
          args := make(map[string]string)
          args["peerIp"] = ipport[0]
          args["peerPort"] = ipport[1]
          args["peerPubKeyHash"] = pubkeyhash

          action := m.PeerRouteActionList[PEER_REQUEST_RESPONSE]
          action(args)
        }
      }
    }
  }
}

/*
	Create a beacon to announce this node's pkh
*/
func (m *MDHT) SelfAnnouncementBeacon(stop *chan bool, initDelaySeconds int, delaySeconds int) error {
  keyManager := mcrypto_key.KeyManagerInstance()
  return m.AnnouncementBeacon(stop, keyManager.GetBasePublicKey(), initDelaySeconds, delaySeconds)
}

/*
	Create a beacon to announce a pkh
*/
func (m *MDHT) AnnouncementBeacon(stop *chan bool, pubKey *rsa.PublicKey, initDelaySeconds int, delaySeconds int) error {
  infoHash, err := mutil.GetInfohashByPubKey(pubKey)
  if err != nil {
    mlog.GetLogger().Error("announcementBeacon - GetInfohashByPubKey failed", err)
    return err
  }

  mlog.GetLogger().Debug(fmt.Sprintf("Starting AnnouncementBeacon - announcing %s", infoHash))

  time.Sleep(time.Duration(initDelaySeconds) * time.Second)
  for {
    select {
    case <-*stop:
      return nil
    default:
      m.DHT.PeersRequest(infoHash, true)
      time.Sleep(time.Duration(delaySeconds) * time.Second)
    }
  }
  return nil
}

/*
	Create a request on the beacon for a specified pkh
*/
func (m *MDHT) RequestBeacon(stop *chan bool, peerPubKeyHash string, action string, actionArgs []string, numBeacon int, intervalBeacon int) {
  mlog.GetLogger().Debug("RequestBeacon - started requesting for peer pubkeyhash", peerPubKeyHash)

  if intervalBeacon <= 0 {
    intervalBeacon = 1
  }

  for {
    select {
    case <-*stop:
      break
    default:
      m.DHT.PeersRequest(peerPubKeyHash, false)
      time.Sleep(time.Duration(intervalBeacon) * time.Second)
    }
  }
}
