package mnet_core_manager

import (
  "../../../runtime"
  "../../../crypto/key"
  "../../if"
  "../../ip"
  "../../security"
  "../../vars"
  "../base"
  "../transport"
  "bytes"
  "fmt"
  "github.com/pkg/errors"
  "gitlab.neji.vm.tc/marconi/log"
  "net"
  "os"
  "strconv"
)

/*
	Encapsulates a connection between two marconi nodes,
*/
type Connection struct {
  ID               uint16
  Port             string
  PhysConnection   net.Conn
  TapConnection    *mnet_if.Interface
  ChannelListen    chan string
  ChannelTransmit  chan string
  ChannelKeyMutate chan string
}

/*
	Return the connection for a specific peer
	This can be extended later to accept a connection ID to support multiple connections per peer, eg for redundancy
*/
func (nm *NetCoreManager) GetConnectionForPeer(peerPubKeyHash string) (*Connection, error) {
  connection, exists := nm.peerConnections[peerPubKeyHash]
  if !exists {
    return nil, errors.New(fmt.Sprintf("No connection found for peer: %s", peerPubKeyHash))
  }
  return connection, nil
}

/*
	Create a new connection object for a peer and return it
*/
func (nm *NetCoreManager) createNewConnectionForPeer(peerPubKeyHash string) *Connection {
  connection := &Connection{}
  connection.ChannelListen = make(chan string)
  connection.ChannelTransmit = make(chan string)
  connection.ChannelKeyMutate = make(chan string)
  nm.peerConnections[peerPubKeyHash] = connection
  return connection
}

/*
	Close a connection for a peer
	If GetConnectionForPeer was extended for multiple connections per peer, this will need to be updated
*/
func (nm *NetCoreManager) closeConnectionForPeer(peerPubKeyHash string) error {
  connection, err := nm.GetConnectionForPeer(peerPubKeyHash)
  if err != nil {
    return err
  }

  // Pushing to these channels terminates the goroutines running for this peer connection
  connection.ChannelTransmit <- "quit"
  connection.ChannelListen <- "quit"
  connection.ChannelKeyMutate <- "quit"

  // close the peer's Tap connection
  connection.TapConnection.Close()
  // close the peer's UDP connection
  err = connection.PhysConnection.Close()
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to close UDP Connection for peer %s with error: %s", peerPubKeyHash, err))
  }

  delete(nm.peerConnections, peerPubKeyHash)
  nm.deallocateConnectionId(connection.ID)

  // Cleanup the pipe states
  nm.pipeStates.Lock()
  if _, exists := (*nm.pipeStates.StateMap)[connection.Port]; exists {
    delete(*nm.pipeStates.StateMap, connection.Port)
  }
  nm.pipeStates.Unlock()

  mlog.GetLogger().Debug("UDP Connection closed for peer", peerPubKeyHash)
  return nil
}

func (nm *NetCoreManager) CheckConnectionPortExists(port string) (bool, uint16) {
  for _, conn := range nm.peerConnections {
    if conn.Port == port {
      return true, conn.ID
    }
  }
  return false, 0
}

/*
	Create a connection to a peer using a TAP driver
*/
func (nm *NetCoreManager) CreateTapConnection(netType NetworkType, netID string, transport mnet_core_transport.Transport, args *mnet_vars.ConnectionArgs) (*Connection, error) {
  var key []byte

  key, err := mcrypto_key.Keyfile_read(args.L2KeyFile)
  /* If the error is file does not exist */
  if err != nil && os.IsNotExist(err) {
    mnet_core_base.Log.Fatalf("Error reading key file: %s", err)
  }

  var primaryDataKey *[]byte
  if args.DataKey != nil {
    b := args.DataKey.Bytes()
    primaryDataKey = &b
  }

  // Allocate a connection id
  connectionId, err := nm.allocateNextConnectionId()
  if err != nil {
    return nil, errors.New("Could not allocate a connection ID")
  }
  connectionIDStr := strconv.Itoa(int(connectionId))

  // Print connection args for debug purposes
  mlog.GetLogger().Debug("Creating tap connection")
  args.DebugPrint()

  //  Create a tap connection
  nm.tapCreationMutex.Lock()
  tapConn := new(mnet_if.Interface)
  tapConn.SetTap(true)
  err = tapConn.OpenTap(mnet_vars.TAP_MTU, connectionIDStr)
  if err != nil {
    mnet_core_base.Log.Fatalf("Error opening a tap device: %s - %s", connectionIDStr, err)
  }
  nm.tapCreationMutex.Unlock()

  mnet_core_base.Log.Debugf("Created tunnel at interface %s with MTU %d", tapConn.GetName(), mnet_vars.TAP_MTU)

  // create a new connection obj for this new connection
  connection := nm.createNewConnectionForPeer(args.PeerPubKeyHash)
  connection.ID = connectionId
  connection.Port = args.LocalPort
  connection.TapConnection = tapConn

  // Add the connection to the bridge
  bridgeInfo, err := nm.GetBridgeInfoForNetwork(netType, netID)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get bridge info: %s", err))
  }
  err = nm.AddConnectionToBridge(bridgeInfo, connectionIDStr, args.RemoteIpAddr)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to add connection to bridge: %s", err))
  }

  // Start a goroutine for mutating key encryption
  go mnet_security.KeepUpdateDataKey(primaryDataKey, args.DataKeySignal, connection.ChannelKeyMutate)

  connection.PhysConnection, err = transport.ListenAndTransmit(
    mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr(), args.LocalPort,
    args.RemoteIpAddr, args.RemotePort, tapConn,
    key, primaryDataKey, args.EncPayload, false,
    connection.ChannelListen, connection.ChannelTransmit)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("CreateTapConnection failed: %s", err))
  }

  return connection, nil
}

func (nm *NetCoreManager) CreateTun(transport mnet_core_transport.Transport, inSignals map[string]chan bool, params ...string) {
  //args into list
  paramList := make([]string, len(params))
  for i, s := range params {
    paramList[i] = s
  }

  mlog.GetLogger().Info("CreateTun: Create Tun Start.. - ", paramList)
  //all arguments
  //l2 or base data integrity check by sha1 on payload
  privateKeyFilePathDataLink := paramList[0]
  //actual encryption key for payload
  payloadEncKey := paramList[1]
  securePayloadEncKey := paramList[2]

  // tap tun num get from allocation
  taptunNumInt, err := nm.allocateNextConnectionId()
  if err != nil {
    // return some error
    return
  }
  taptunNum := strconv.Itoa(int(taptunNumInt))

  mLocalIpAddr := paramList[4]
  mLocalIpMask := paramList[5]
  mLocalIpGatewayIpAddr := paramList[6]
  port := paramList[7]
  targetIpAddr := ""
  targetPort := ""
  if len(paramList) > 8 {
    targetIpAddr = paramList[8]
    targetPort = paramList[9]
  }

  var key []byte
  //var currentBlockTimeIndex uint16

  /*  key, err := keyfile_read(l2EncKey)
   */ /* If the error is file does not exist */ /*
  	if err != nil && os.IsNotExist(err) {
  */ /* Auto-generate the key file */ /*
  		key, err = keyfile_generate(l2EncKey)
  		if err != nil {
  			log.Fatalf("Error generating key file: %s\n", err)
  		}
  	} else if err != nil {
  		log.Fatalf("Error reading key file: %s\n", err)
  	}*/
  key = mcrypto_key.GetLinkLayerEncryptionKey(privateKeyFilePathDataLink, true)

  //data key for encrypt/decript l2/tap packets (Ethernet)
  dataKey, err := mcrypto_key.Keyfile_read(payloadEncKey)
  /* If the error is file does not exist */
  if err != nil && os.IsNotExist(err) {
    mnet_core_base.Log.Fatalf("Error data key file invalid: %s\n", err)
  }
  mlog.GetLogger().Info("dataKey: ", dataKey)
  isSecure := true
  if securePayloadEncKey == "nosecure" {
    isSecure = false
  }

  nm.tunCreationMutex.Lock()
  mlog.GetLogger().Info("CreateTapConnection: Creating Tap EdgeNodeFinderInstance - #", taptunNum)

  //TODO: based on os
  //TODO: redo
  /* Create a tap interface */
  tunConn := new(mnet_if.Interface)
  //mtu := make(uint)
  //mtu = TUN_MTU //TAP_MTU
  //err = tap_conn.Open(TUN_MTU, paramList[4], "tap")
  //err = tap_conn.OpenTap(TAP_MTU, taptunNum)
  err = tunConn.OpenTun(mnet_vars.TUN_MTU, taptunNum)
  if err != nil {
    mnet_core_base.Log.Fatalf("Error opening a tun device: %s - %s\n", taptunNum, err)
  }
  nm.tunCreationMutex.Unlock()

  mnet_core_base.Log.Infof("Created tunnel at interface %s with MTU %d\n\n", tunConn.GetName(), mnet_vars.TAP_MTU)
  mnet_core_base.Log.Info("Starting marconitunnel/mtun...")
  primaryDataKey := make([]byte, mnet_vars.DATA_KEY_SIZE)
  primaryDataKey = dataKey

  //updating payload aes encryption key
  // TODO: ayuen we need to make sure this is hooked up to dh key info key signal
  signal := make(chan *bytes.Buffer)
  go mnet_security.KeepUpdateDataKey(&primaryDataKey, &signal, nil)

  //NOTE: disabling for now
  //go StartSocket5Proxy()
  //go StartSocket5Proxy2()

  //TODO: do for tun/l3 pointopoint = client mode
  //use mLocalIpGatewayIpAddras peer IP instead of targetIpAddr since it's public or internet world which
  //marconi does not have direct access
  // TODO: recursively inline this function later so that NetCoreManager can keep track of the tuns created
  mnet_ip.ConfigMconnIpAddress(taptunNum, mLocalIpAddr, mLocalIpMask, mLocalIpGatewayIpAddr, targetIpAddr)

  transport.ListenAndTransmit("0.0.0.0", port, targetIpAddr, targetPort, tunConn, key, &primaryDataKey, isSecure, true, nil, nil)
}
