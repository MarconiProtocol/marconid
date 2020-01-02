package magent_edge_route_request

import (
  "../../../core/crypto/dh"
  "../../../core/crypto/key"
  "../../../core/net/core/manager"
  "../../../core/net/core/transport/udp"
  "../../../core/net/vars"
  "../../../core/peer"
  "../../service/edge_route"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "net"
  "strconv"
  "sync"
)

const (
  EDGE_PORT_BASE uint16 = 30000
  EDGE_PORT_MAX  uint16 = 31024
)

type EdgeRouteRequestManager struct {
  allocatedPorts      []bool
  portAllocationMutex sync.Mutex
}

var edgeRouteRequestManager *EdgeRouteRequestManager
var edgeRouteRequestManagerOnce sync.Once

func Instance() *EdgeRouteRequestManager {
  edgeRouteRequestManagerOnce.Do(func() {
    edgeRouteRequestManager = &EdgeRouteRequestManager{}
    edgeRouteRequestManager.initialize()
  })
  return edgeRouteRequestManager
}

func (e *EdgeRouteRequestManager) initialize() {
  e.allocatedPorts = make([]bool, EDGE_PORT_MAX-EDGE_PORT_BASE+1)
}

/*
	Request an edge route to a peer
*/
func (e *EdgeRouteRequestManager) RequestRoute(remoteNodeIpAddr string, pkhIdentifier string) error {
  peerManager := mpeer.PeerManagerInstance()

  // Send RPC to remote node, requesting for an edge route to be reserved for the client
  edgeRouteResp, err := sendRequestEdgeRouteRPC(remoteNodeIpAddr, pkhIdentifier)
  if err != nil {
    mlog.GetLogger().Errorf("error during sendRequestEdgeRouteRPC %s", err.Error())
    return err
  } else {
    mlog.GetLogger().Info(fmt.Sprintf("Received edgeRouteResp from RequestEdgeRoute, node with pkh %s says we can connect on port %s using host %s",
      edgeRouteResp.edgeNodePubKeyHash, edgeRouteResp.port, edgeRouteResp.host,
    ))

    // TODO : this flow may change when we integrate the SC piece
    // The remote node reserved an edge route for this client, add it as an edge peer
    peerManager.AddPeer(mpeer.EDGE_PEER, edgeRouteResp.edgeNodePubKeyHash)

    // Initiate a public key exchange with the target peer (when necessary)
    err := mcrypto_key.KeyManagerInstance().InitiatePublicKeyExchange(remoteNodeIpAddr, edgeRouteResp.edgeNodePubKeyHash)
    if err != nil {
      mlog.GetLogger().Errorf("Failed to perform pk exchange with edge route provider: %s", err)
      return err
    }

    // Initiate a Diffie-Hellman key exchange with the target peer (when necessary)
    err = mcrypto_dh.DHExchangeManagerInstance().InitiateDHKeyExchange(edgeRouteResp.edgeNodePubKeyHash)
    if err != nil {
      mlog.GetLogger().Errorf("Failed to perform dh exchange with edge route provider: %s", err)
      return err
    }

    // Create edge route to the route provider node
    err = e.startRouteToNode(edgeRouteResp.edgeNodePubKeyHash, remoteNodeIpAddr, edgeRouteResp.port, edgeRouteResp.host, pkhIdentifier)
    if err != nil {
      mlog.GetLogger().Fatal("Failed to establish route to node")
      return err
    }
  }
  return nil
}

func (e *EdgeRouteRequestManager) startRouteToNode(peerPubKeyHash string, remoteIpAddr string, port string, assignedHost string, identifier string) error {

  mlog.GetLogger().Infof("pkh: %s remoteIpAddr %s port %s assigned host %s identifier %s", peerPubKeyHash, remoteIpAddr, port, assignedHost, identifier)

  // get DH key info for peer
  dhKeyInfo, err := mcrypto_dh.DHExchangeManagerInstance().GetDHKeyInfo(peerPubKeyHash)
  if err != nil {
    mlog.GetLogger().Fatal("Could not find exchanged DH info for peer: ", peerPubKeyHash)
    return err
  }

  // allocate an unused port
  ourPort, err := e.allocateNextEdgePort()
  if err != nil {
    //return an error instead
    mlog.GetLogger().Error("Could not allocated an edge port")
    return err
  }

  // create a bridge to bridge all mpipes
  ncm := mnet_core_manager.GetNetCoreManager()
  bridgeInfo, newlyAllocated, err := ncm.GetOrAllocateBridgeInfoForNetwork(mnet_core_manager.EDGE_NET, identifier)
  if err != nil {
    mlog.GetLogger().Info(fmt.Sprintf("Could not get or allocate a bridge: %s", err.Error()))
    return err
  } else {
    if newlyAllocated {
      ncm.CreateBridge(bridgeInfo, assignedHost, magent_edge_route_provider.EDGE_NETWORK_NETMASK, false)
    }
    ncm.AssignIpAddrToBridge(bridgeInfo, assignedHost, magent_edge_route_provider.EDGE_NETWORK_NETMASK)
  }

  // use placeholders for a bunch of this stuff for now
  args := mnet_vars.ConnectionArgs{
    L2KeyFile:      "/opt/marconi/etc/marconid/l2.key",
    EncPayload:     true,
    LocalPort:      strconv.Itoa(int(ourPort)),
    RemoteIpAddr:   remoteIpAddr,
    RemotePort:     port,
    DataKey:        dhKeyInfo.SymmetricKeyBytes,
    DataKeySignal:  dhKeyInfo.KeySignal,
    PeerPubKeyHash: peerPubKeyHash,
  }

  // Transforming MpipeArgs to TapArgs, but I dont think there really is a difference between the two...
  mlog.GetLogger().Info(fmt.Sprintf("==> Creating EdgeMPipe to targetIP: %s, targetPort: %s", args.RemoteIpAddr, args.RemotePort))
  conn, err := ncm.CreateTapConnection(mnet_core_udp.GetUDPTransport(), &args)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("error: %s \n", err.Error()))
    return err
  }

  // Add the connection to the bridge
  err = ncm.AddConnectionToBridge(bridgeInfo, strconv.Itoa(int(conn.ID)), "")
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed to add connection to bridge: %s", err.Error()))
    return err
  }
  return nil
}

func (e *EdgeRouteRequestManager) allocateNextEdgePort() (uint16, error) {
  e.portAllocationMutex.Lock()
  defer e.portAllocationMutex.Unlock()

  for i, allocated := range e.allocatedPorts {
    if !allocated {
      e.allocatedPorts[i] = true
      // check if the port is available
      portCandidate := uint16(i) + EDGE_PORT_BASE
      addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(int(portCandidate)))
      udpConn, err := net.ListenUDP("udp", addr)
      // If the port is good, we can return it
      if err == nil {
        _ = udpConn.Close()
        return portCandidate, nil
      }
    }
  }
  return 0, errors.New("edge_route::allocateNextEdgePort -> couldn't find an unallocated edge port")
}
