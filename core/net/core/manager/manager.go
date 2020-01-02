package mnet_core_manager

import (
  "../../../peer"
  "../../vars"
  "../base"
  "../transport/udp"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "strconv"
  "sync"
)

type NetworkType int

const (
  SERVICE_NET NetworkType = iota
  EDGE_NET
)

type NetCoreManager struct {
  netIDtoBridgeInfoMap map[string]*BridgeInfo  // Maps Net ID to bridge information
  peerConnections      map[string]*Connection  // Maps peer pubkeyhash to connection obj
  ifaceIDAllocations   *InterfaceIDAllocations // Keeps track of all allocated iface IDs
  pipeStates           *PipeStateMap           // Keeps track of mpipe states

  tapCreationMutex sync.Mutex
  tunCreationMutex sync.Mutex
  sync.Mutex
}

var once sync.Once
var netCoreManager *NetCoreManager

func GetNetCoreManager() *NetCoreManager {
  once.Do(func() {
    netCoreManager = createNetCoreManager()
  })
  return netCoreManager
}

func createNetCoreManager() *NetCoreManager {
  var err error
  mnet_core_base.Log, err = mlog.GetLogInstance("netcore")
  if err != nil {
    mlog.GetLogger().Error("failed to init: log - netcore", err)
  }

  return &NetCoreManager{
    netIDtoBridgeInfoMap: make(map[string]*BridgeInfo),
    ifaceIDAllocations:   initializeInterfaceIdAllocations(),
    pipeStates:           NewPipeStates(),
    peerConnections:      make(map[string]*Connection),
  }
}

/*
	Create a MPipe to a peer, with the given arguments
*/
func (nm *NetCoreManager) CreateMPipe(args *mnet_vars.ConnectionArgs) error {
  // we need to make sure we check the pipe state at this point
  nm.pipeStates.Lock()
  currentState := (*nm.pipeStates.StateMap)[args.LocalPort]
  if currentState == SUCCESS || currentState == ATTEMPTING {
    // no-op since we are done... OR trying to create the pipe through another goroutine
    mlog.GetLogger().Warn("CreateConnectionToServicePeer skipping MPipe creation due to pipe's currentState", currentState)
    nm.pipeStates.Unlock()
    return nil
  } else {
    // set the state to ATTEMPTING to signify that there is already a callback that is trying to create an mpipe
    (*nm.pipeStates.StateMap)[args.LocalPort] = ATTEMPTING
  }
  nm.pipeStates.Unlock()
  status := SUCCESS

  // Create a tap connection and attach it to the bridge with netType SERVICE_NET and netId 'main'
  conn, err := nm.CreateTapConnection(mnet_core_udp.GetUDPTransport(), args)
  if err != nil {
    status = UNATTEMPTED
    return errors.New(fmt.Sprintf("Failed to create tap connection: %s", err))
  }

  // Add the connection to the bridge
  bridgeInfo, err := nm.GetBridgeInfoForNetwork(SERVICE_NET, "main")
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to get bridge info: %s", err))
  }
  err = nm.AddConnectionToBridge(bridgeInfo, strconv.Itoa(int(conn.ID)), args.RemoteIpAddr)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to add connection to bridge: %s", err))
  }

  // If we are at this point we have started the transmit and listen goroutines via MPipe() and we can claim we have finished our attempt at creating a pipe
  nm.pipeStates.Lock()
  (*nm.pipeStates.StateMap)[args.LocalPort] = status
  nm.pipeStates.Unlock()

  return nil
}

/*
	Close a MPipe to a peer
*/
func (nm *NetCoreManager) CloseMPipe(peer *mpeer.Peer) {
  // Close the connection for the peer
  err := nm.closeConnectionForPeer(peer.PubKeyHash)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed to close a connection for peer %s, error: %s", peer.PubKeyHash, err))
  }
}
