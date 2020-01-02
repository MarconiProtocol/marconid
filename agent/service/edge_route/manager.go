package magent_edge_route_provider

import (
  "../../../core/crypto/dh"
  "../../../core/net/core/manager"
  "../../../core/net/core/transport/udp"
  "../../../core/net/vars"
  "errors"
  "fmt"
  "github.com/MarconiProtocol/log"
  "net"
  "strconv"
  "sync"
  "time"
)

const (
  EDGE_PORT_BASE uint16 = 50000
  EDGE_PORT_MAX  uint16 = 51000
  NETWORK_MIN    uint8  = 11
  NETWORK_MAX    uint8  = 254

  // TEMP
  EDGE_NETWORK_FIRST_OCTET   = "172."
  EDGE_PROVIDER_NETWORK_HOST = ".10.1"
  EDGE_CLIENT_NETWORK_HOST   = ".10.2"
  EDGE_NETWORK_NETMASK       = 30
)

type EdgeRouteManager struct {
  edgeConnections map[string]*EdgeConnection

  unallocatedPorts    []uint16 // Ports that are usable and unallocated, initially assume every port is usable between PORT_BASE AND PORT_MAX
  allocatedPorts      []uint16 // Ports that are currently allocated to an edge connection
  unusablePorts       []uint16 // Ports that have been tested and were not available, categorize them into a separate bucket to avoid frequent retests
  portAllocationMutex sync.Mutex

  allocatedNetworks      []bool
  networkAllocationMutex sync.Mutex
}

type EdgeConnection struct {
  //PubKeyHash string
  AssignedIPNetwork uint8 // Assigned 2nd octet of the network portion of IP address, eg 172.<AssignedIPNetwork>.10.10
  Port              uint16

  //Connected bool	// Need some kind of reservation concept
}

var instance *EdgeRouteManager
var once sync.Once

func EdgeConnectionManagerInstance() *EdgeRouteManager {
  once.Do(func() {
    instance = &EdgeRouteManager{}
    instance.initialize()
  })
  return instance
}

func (e *EdgeRouteManager) initialize() {
  e.portAllocationMutex.Lock()
  defer e.portAllocationMutex.Unlock()

  e.unallocatedPorts = []uint16{}
  e.allocatedPorts = []uint16{}
  e.unusablePorts = []uint16{}
  for port := EDGE_PORT_BASE; port < EDGE_PORT_MAX; port++ {
    e.unallocatedPorts = append(e.unallocatedPorts, port)
  }

  e.allocatedNetworks = make([]bool, NETWORK_MAX-NETWORK_MIN+1)

  e.edgeConnections = make(map[string]*EdgeConnection)
}

func (e *EdgeRouteManager) RequestEdgeConnection(clientIp string) (*EdgeConnection, error) {
  // TODO: This function should have more args to identify the requested edge connection,
  //  since a single FEO can request an edge connection on behalf of multiple end-users
  // Check if there is an existing edge connection to return
  if edgeConnection, exists := e.edgeConnections[clientIp]; exists {
    return edgeConnection, nil
  }

  // Get port and network allocations
  port, err := e.allocatePort()
  if err != nil {
    return nil, err
  }

  network, err := e.allocateIPNetwork()
  if err != nil {
    return nil, err
  }

  edgeConnection := EdgeConnection{
    network,
    port,
  }

  e.edgeConnections[clientIp] = &edgeConnection

  return &edgeConnection, nil
}

func StartEdgeRoute(peerPubKeyHash string, port uint16, identifier string, assignedNetwork string) {
  // TODO : do a dumb thing for now
  //  basically waits till we have completed the DH exchange...
  var dhKeyInfo *mcrypto_dh.DHKeyInfo
  var err error
  for {
    dhKeyInfo, err = mcrypto_dh.DHExchangeManagerInstance().GetDHKeyInfo(peerPubKeyHash)
    if err != nil {
      mlog.GetLogger().Debug("Could not find exchanged DH info for peer: WAITING MORE", peerPubKeyHash)
    } else {
      break
    }
    time.Sleep(time.Duration(time.Second * 2))
  }

  // placeholders for now, easier to change later
  args := mnet_vars.ConnectionArgs{
    L2KeyFile:      "/opt/marconi/etc/marconid/l2.key",
    EncPayload:     true,
    LocalPort:      strconv.Itoa(int(port)),
    RemoteIpAddr:   "",
    RemotePort:     "",
    DataKey:        dhKeyInfo.SymmetricKeyBytes,
    DataKeySignal:  dhKeyInfo.KeySignal,
    PeerPubKeyHash: peerPubKeyHash,
  }

  mlog.GetLogger().Info(fmt.Sprintf("==> Accepting an edge route with identity %s to our port: %d", identifier, port))

  // TODO: ayuen -> probably need to do more

  ncm := mnet_core_manager.GetNetCoreManager()
  conn, err := ncm.CreateTapConnection(mnet_core_udp.GetUDPTransport(), &args)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("error: %s \n", err.Error()))
  }

  // Add the connection to the bridge
  // create a bridge to bridge all mpipes
  bridgeInfo, newlyAllocated, err := ncm.GetOrAllocateBridgeInfoForNetwork(mnet_core_manager.EDGE_NET, identifier)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Could not get or allocate a bridge: %s", err.Error()))
  } else {
    mlog.GetLogger().Info(fmt.Sprintf("Bridge info %s is newlyAllocated %b", bridgeInfo, newlyAllocated))
    ipAddr := EDGE_NETWORK_FIRST_OCTET + assignedNetwork + EDGE_PROVIDER_NETWORK_HOST
    if newlyAllocated {
      ncm.CreateBridge(bridgeInfo, ipAddr, EDGE_NETWORK_NETMASK, false)
    }
    ncm.AssignIpAddrToBridge(bridgeInfo, ipAddr, EDGE_NETWORK_NETMASK)
  }
  err = ncm.AddConnectionToBridge(bridgeInfo, strconv.Itoa(int(conn.ID)), "")
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed to add connection to bridge: %s", err.Error()))
  }

}

func (e *EdgeRouteManager) allocatePort() (uint16, error) {
  // make sure this function is locked
  e.portAllocationMutex.Lock()
  defer e.portAllocationMutex.Unlock()

  // check if the port is available
  var freePortCandidate uint16
  for {
    // first check if there are any unallocated ports
    if len(e.unallocatedPorts) == 0 {
      return 0, errors.New("EdgeRouteManager::allocatePort - ran out of ports to allocate")
    }

    // simply choose the first unallocated port as the free port candidate
    freePortCandidate = e.unallocatedPorts[0]
    e.unallocatedPorts = e.unallocatedPorts[1:]

    addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(int(freePortCandidate)))
    udpConn, err := net.ListenUDP("udp", addr)
    // If there is an error, consider that port to be unusable
    if err != nil {
      e.unusablePorts = append(e.unusablePorts, freePortCandidate)
    } else {
      e.allocatedPorts = append(e.allocatedPorts, freePortCandidate)
      _ = udpConn.Close()
      break
    }
  }
  return freePortCandidate, nil
}

func (e *EdgeRouteManager) deallocatePort(port uint16) {
  e.portAllocationMutex.Lock()
  defer e.portAllocationMutex.Unlock()

  //TODO dumb lookup
}

// KISS FOR NOW since this will need to be replaced by calls to SC most likely
func (e *EdgeRouteManager) allocateIPNetwork() (uint8, error) {
  e.networkAllocationMutex.Lock()
  defer e.networkAllocationMutex.Unlock()

  numNetworks := NETWORK_MAX - NETWORK_MIN + 1
  var idx uint8
  for idx = 0; idx < numNetworks; idx++ {
    if e.allocatedNetworks[idx] == false {
      e.allocatedNetworks[idx] = true
      return idx + NETWORK_MIN, nil
    }
  }
  return 0, errors.New("edge_route::allocateIPNetwork -> couldnt find an unallocated network")
}
