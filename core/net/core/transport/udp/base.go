package mnet_core_udp

import (
  "../../../../../core"
  "../../../../config"
  "../../../../runtime"
  "../../../../../util"
  "../../../if"
  "../../../packet/filter"
  "../../../vars"
  "../../base"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "errors"
  "fmt"
  "net"
  "strings"
  "sync"
  "time"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"

  "git.marconi.org/marconiprotocol/sdk/packet/filter"
)

type UDPTransport struct{}

var once sync.Once
var udpTransport *UDPTransport

func GetUDPTransport() *UDPTransport {
  once.Do(func() {
    udpTransport = &UDPTransport{}
  })

  return udpTransport
}

func (udpt *UDPTransport) ListenAndTransmit(
  localIpAddr string, localPort string,
  remoteIpAddr string, remotePort string,
  tapConn *mnet_if.Interface,
  key []byte, dataKey *[]byte,
  isSecure bool, isTun bool,
  listenSignalChannel chan string, transmitSignalChannel chan string) (net.Conn, error) {

  /* Parse & resolve the peer address, if it was provided */
  var peerIpAddr *net.UDPAddr
  var peerDiscoveryChannel chan net.UDPAddr

  //Parse & resolve local address
  localUDPAddr, err := net.ResolveUDPAddr("udp", localIpAddr+":"+localPort)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to resolve local addr %s:%s - err: %s", localIpAddr, localPort, err))
  }

  //Create a UDP physical connection
  physConn, err := net.ListenUDP("udp", localUDPAddr)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to create UDP socket %s - err: %s", localUDPAddr, err))
  }

  if remoteIpAddr != "" {
    remoteTarget := remoteIpAddr + ":" + remotePort
    peerIpAddr, err = net.ResolveUDPAddr("udp", remoteTarget)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Failed to resolve remote addr %s:%s - err: %s", remoteIpAddr, remotePort, err))
    }
    peerDiscoveryChannel = nil
  } else {
    //Otherwise, prepare a channel that the listen() goroutine will forward discovered peers through
    peerIpAddr = nil
    peerDiscoveryChannel = make(chan net.UDPAddr)
  }

  go udpt.Listen(physConn, tapConn, peerIpAddr, key, dataKey, isSecure, peerDiscoveryChannel, isTun, listenSignalChannel)
  go udpt.Transmit(physConn, tapConn, peerIpAddr, key, dataKey, isSecure, peerDiscoveryChannel, isTun, transmitSignalChannel)

  return physConn, nil
}





//To forward data from physical interface to tap/tun/virtual link interface
func (udpt *UDPTransport) Listen(physConn *net.UDPConn, tapConn *mnet_if.Interface, peerAddr *net.UDPAddr, key []byte,
  dataKey *[]byte, isSecure bool, peerDiscoveryChannel chan net.UDPAddr, isTun bool, signalChannel chan string) {

  packet := make([]byte, mnet_vars.UDP_MTU)
  var decapsulatedFrames []byte
  var invalid error = nil

  var currentPeerAddr net.UDPAddr
  var peerDiscovery bool

  hmacH := hmac.New(sha256.New, key)

  /* If a peer was specified, fill in our peer information */
  if peerAddr != nil {
    currentPeerAddr.IP = peerAddr.IP
    currentPeerAddr.Port = peerAddr.Port
    mnet_core_base.Log.Debugf("Starting udp->tap forwarding with peer %s:%d...", currentPeerAddr.IP, currentPeerAddr.Port)
    peerDiscovery = false
  } else {
    mnet_core_base.Log.Debug("Starting udp->tap forwarding with peer discovery...")
    peerDiscovery = true
  }

  var crashCounterInterface uint = 0

  filterManager := mpacket_filter.GetFilterManagerInstance()
  resultChannel := make(chan mpacket_filter.PacketFilterResult)

  for {
    select {
    case msg := <-signalChannel:
      if strings.Compare(msg, "quit") == 0 {
        return
      }
    default:
      /* Receive an encapsulated frame packet through UDP */
      physConn.SetDeadline(time.Now().Add(mnet_core_base.SOCKET_TIMEOUT_SECONDS * time.Second))
      n, raddr, err := physConn.ReadFromUDP(packet)
      if err != nil {
        mnet_core_base.Log.Warnf("Error reading from UDP socket: %s", err)
        continue
      }

      /* If peer discovery is off, ensure the received packge is from our specified peer */
      if !peerDiscovery && (!raddr.IP.Equal(currentPeerAddr.IP) || raddr.Port != currentPeerAddr.Port) {
        continue
      }

      mnet_core_base.Log.Debug("<- udp  | Encapsulated frame:")
      mnet_core_base.Log.Debugf("        | from Peer %s:%d", raddr.IP, raddr.Port)
      packetDump := strings.Split(hex.Dump(packet[0:n]), "\n")
      for i := 0; i < len(packetDump)-1; i++ {
        mnet_core_base.Log.Debug(packetDump[i])
      }

      /* Decapsulate the frame, skip it if it's invalid */
      decapsulatedFrames, invalid = mnet_core_base.DecodeFrame(packet[0:n], hmacH)
      if invalid != nil {
        mnet_core_base.Log.Debugf("<- udp  | Frame discarded! Size: %d, Reason: %s", n, invalid.Error())
        mnet_core_base.Log.Debugf("        | from Peer %s:%d", raddr.IP, raddr.Port)
        mnet_core_base.Log.Debug(hex.Dump(packet[0:n]))
        continue
      }

      decapsulatedFramesDecoded := decapsulatedFrames
      if isSecure {
        var crypter mutil.SymmetricCrypter = mutil.NewAesCrypter(*dataKey)
        tmpFrame, err := crypter.Decrypt(decapsulatedFrames)

        decapsulatedFramesDecoded = make([]byte, len(tmpFrame))
        decapsulatedFramesDecoded = tmpFrame

        if err != nil {
          mnet_core_base.Log.Fatal("Error: decrypted:", len(decapsulatedFramesDecoded), "org/encrypted:", len(decapsulatedFrames), "n", n, "err:", err)
        } else {
          mnet_core_base.Log.Debug("decrypted:", len(decapsulatedFramesDecoded), "org/encrypted:", len(decapsulatedFrames), "after l2 checksum strip:", decapsulatedFrames, "n:", n)
        }
      } else {
        mnet_core_base.Log.Debug("decrypt: no encryption enabled")
      }

      /* If peer discovery is on and the peer is new, save the discovered peer */
      if peerDiscovery && (!raddr.IP.Equal(currentPeerAddr.IP) || raddr.Port != currentPeerAddr.Port) {
        currentPeerAddr.IP = raddr.IP
        currentPeerAddr.Port = raddr.Port
        /* Send the new peer info to our transmit() goroutine via channel */
        peerDiscoveryChannel <- currentPeerAddr

        mnet_core_base.Log.Debugf("Discovered peer %s:%d!", currentPeerAddr.IP, currentPeerAddr.Port)
      }

      mnet_core_base.Log.Debug("-> tap  | Decapsulated frame from peer:")
      mnet_core_base.Log.Debug(hex.Dump(decapsulatedFramesDecoded))
      mnet_core_base.Log.Debug("decoded frame:", decapsulatedFramesDecoded, len(decapsulatedFramesDecoded))

      // only ask filterManager to process packets if packet filtering is enabled
      if mconfig.GetAppConfig().Global.PacketFiltersEnabled {
        // Marshall bytes into gopacket
        ethPacket := gopacket.NewPacket(decapsulatedFramesDecoded, layers.LayerTypeEthernet, gopacket.Default)

        // Process the packet
        go filterManager.ProcessPacket(&ethPacket, &resultChannel)

        // Wait for the signal from filter manager packet signal channel
        result := <-resultChannel
        switch result.Action {
        case m_packet_filter.ACCEPT:
          // No-op
        case m_packet_filter.DROP:
          mnet_core_base.Log.Debug("Dropped packet - ", result.FilterResponse.Msg)
          // Drop the packet by not forwarding it on to virtual device
          continue
        }
      }

      if isTun {
        //NOTE: darwin/osx does not need to go though clean up path since it requires ot have extra 4 bytes on header
        if mruntime.GetMRuntime().GetRuntimeOS() != mcore.TYPE_OS_DARWIN {
          //for now try slice and repack
          adjustedFrame := decapsulatedFramesDecoded[4:]    // bsd link layer adjustment
          decapsulatedFramesDecoded = make([]byte, len(adjustedFrame))
          decapsulatedFramesDecoded = adjustedFrame
        } else {
          BSDLinkLayerBytes := []byte{0, 0, 0, 2}
          possibleBSDLinkLayerBytes := decapsulatedFramesDecoded[:4]
          _ = BSDLinkLayerBytes
          _ = possibleBSDLinkLayerBytes
          mnet_core_base.Log.Debug("Regular Linux style tun traffic, not BSD style.")
        }
      } else {
        mnet_core_base.Log.Debug("net_core/listen: mode: tap")
      }

      _, err = tapConn.Write(decapsulatedFramesDecoded)
      if err != nil {
        crashCounterInterface++
        mnet_core_base.Log.Debug("interface: file", tapConn.GetFdOS())
        mnet_core_base.Log.Debug("interface: pointer", tapConn.GetFd())
        mnet_core_base.Log.Debug("decoded frame:", decapsulatedFramesDecoded, len(decapsulatedFramesDecoded))
        mnet_core_base.Log.Error("Error writing to tap device:", err)
      }
    }
  }
}








//To forward data from tap/tun/virtual link interface to physical interface
func (udpt *UDPTransport) Transmit(physConn *net.UDPConn, tapConn *mnet_if.Interface, peerAddr *net.UDPAddr, key []byte, dataKey *[]byte,
  isSecure bool, peerDiscoveryChannel chan net.UDPAddr, isTun bool, signalChannel chan string) {

  /* Raw tap frame received */
  //var frame []byte
  //NOTE: this value matter with MTU since it's writing into peer and this buffer is max value
  // but this might not matter with tun interace. try with 1500 - pedding for encyrptions
  frame := make([]byte, mnet_vars.TAP_MTU+14)
  /* Encapsulated frame and error */
  var encapsulatedFrames []byte
  var invalid error = nil
  /* Peer address */
  var currentPeerAddr net.UDPAddr

  var peerDiscovery bool

  //var currentBlockTimeIndex uint16

  /* Initialize our HMAC-SHA256 hash context */
  hmacH := hmac.New(sha256.New, key)

  /* If a peer was specified, fill in our peer information */
  if peerAddr != nil {
    currentPeerAddr.IP = peerAddr.IP
    currentPeerAddr.Port = peerAddr.Port
    peerDiscovery = false
  } else {
    peerDiscovery = true
    /* Otherwise, wait for the listen() goroutine to discover a peer */
    currentPeerAddr = <-peerDiscoveryChannel
  }

  mnet_core_base.Log.Debugf("Starting tap->udp forwarding with peer %s:%d...", currentPeerAddr.IP, currentPeerAddr.Port)

  for {
    select {
    case msg := <-signalChannel:
      if strings.Compare(msg, "quit") == 0 {
        return
      }
    default:
      /* If peer discovery is on, check for any newly discovered peers */
      if peerDiscovery {
        select {
        case currentPeerAddr = <-peerDiscoveryChannel:
          mnet_core_base.Log.Debugf("Starting tap->udp forwarding with peer %s:%d...", currentPeerAddr.IP, currentPeerAddr.Port)
        default:
        }
      }

      /* Read a raw frame from our tap device */
      n, err := tapConn.Read(frame)
      if err != nil {
        mnet_core_base.Log.Debug("Tap/Tun Interface Read: ", n)
        mnet_core_base.Log.Debug("Interface:", tapConn.GetFdOS())
        mnet_core_base.Log.Fatalf("Error reading from tap device: %s", err)
      }

      mnet_core_base.Log.Debug("<- tap  | Plaintext frame to peer:")
      frameDump := strings.Split(hex.Dump(frame[0:n]), "\n")
      for i := 0; i < len(frameDump)-1; i++ {
        mnet_core_base.Log.Debug(frameDump[i])
      }
      mnet_core_base.Log.Debug("raw frame: ", frame, len(frame))

      rawFrame := frame[:n]
      if isTun {
        // we do not need to inject with DARWIN/OSX system since system will add extra 4 bytes
        if mruntime.GetMRuntime().GetRuntimeOS() != mcore.TYPE_OS_DARWIN {
          //for now try slice and repack
          adjustedFrame := append([]byte{0, 0, 0, 2}, rawFrame...)
          rawFrame = make([]byte, len(adjustedFrame))
          rawFrame = adjustedFrame
        }
      } else {
        mnet_core_base.Log.Debug("net_core/transmit: mode: tap")
      }

      //encryption on data
      if isSecure {
        var crypter mutil.SymmetricCrypter = mutil.NewAesCrypter(*dataKey)
        dataFrame, err := crypter.Encrypt(rawFrame)
        /* Encapsulate the frame, skip it if it's invalid */
        encapsulatedFrames, invalid = mnet_core_base.EncodeFrame(dataFrame, hmacH)
        if err != nil {
          mnet_core_base.Log.Fatal("Error: encrypted: ", len(dataFrame), "org:", len(rawFrame), "err:", err)
        } else {
          mnet_core_base.Log.Debug("encrypted: ", len(dataFrame), "org:", len(rawFrame))
        }
      } else {
        mnet_core_base.Log.Debug("encrypt: no encryption enabled")
        /* Encapsulate the frame, skip it if it's invalid */
        encapsulatedFrames, invalid = mnet_core_base.EncodeFrame(rawFrame, hmacH)
      }

      if invalid != nil {
        mnet_core_base.Log.Debugf("-> udp  | Frame discarded! Size: %d, Reason: %s", n, invalid.Error())
        mnet_core_base.Log.Debug(hex.Dump(rawFrame))
        continue
      }

      mnet_core_base.Log.Debug("-> udp  | Encapsulated frame to peer:", currentPeerAddr.IP.String())
      encFrameDump := strings.Split(hex.Dump(encapsulatedFrames), "\n")
      for i := 0; i < len(encFrameDump)-1; i++ {
        mnet_core_base.Log.Debug(encFrameDump[i])
      }

      mnet_core_base.Log.Debug("encapsulated encrypted frame: ", encapsulatedFrames, len(encapsulatedFrames))

      physConn.SetDeadline(time.Now().Add(mnet_core_base.SOCKET_TIMEOUT_SECONDS * time.Second))
      _, err = physConn.WriteToUDP(encapsulatedFrames, &currentPeerAddr)
      if err != nil {
        mnet_core_base.Log.Warnf("Error writing to UDP socket: %s", err)
      }
    }
  }
}
