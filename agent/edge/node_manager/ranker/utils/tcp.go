package magent_edge_ranker_utils

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "github.com/MarconiProtocol/log"
  "math/rand"
  "net"
  "strconv"
  "strings"
  "sync"
  "time"
)

const (
  FIN = 1  // 00 0001
  SYN = 2  // 00 0010
  RST = 4  // 00 0100
  PSH = 8  // 00 1000
  ACK = 16 // 01 0000
  URG = 32 // 10 0000
)

type TCPHeader struct {
  Source      uint16
  Destination uint16
  SeqNum      uint32
  AckNum      uint32
  DataOffset  uint8 // 4 bits
  Reserved    uint8 // 3 bits
  ECN         uint8 // 3 bits
  Ctrl        uint8 // 6 bits
  Window      uint16
  Checksum    uint16 // Kernel will set this if it's 0
  Urgent      uint16
  Options     []TCPOption
}

type TCPOption struct {
  Kind   uint8
  Length uint8
  Data   []byte
}

const DEFAULT_DELAY_BETWEEN_LATENCY_CHECK_TCP = 2000

func TimeTCPSynAck(localAddr string, remoteAddr string, port uint16) time.Duration {
  var wg sync.WaitGroup
  wg.Add(1)
  var receiveTime time.Time

  addrs, err := net.LookupHost(remoteAddr)
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("Error resolving %s. %s\n", remoteAddr, err))
  }
  remoteAddr = addrs[0]

  // Run go-routine to start waiting for a TCP Syn-Ack
  go func() {
    receiveTime = receiveTCPSynAck(localAddr, remoteAddr)
    wg.Done()
  }()
  time.Sleep(1 * time.Millisecond)

  // Send a TCP Syn packet
  sendTime := sendTCPSyn(localAddr, remoteAddr, port)

  // Timeout
  waitTimeout(&wg, time.Millisecond*DEFAULT_DELAY_BETWEEN_LATENCY_CHECK_TCP)

  return receiveTime.Sub(sendTime)
}

func sendTCPSyn(laddr string, raddr string, port uint16) time.Time {
  packet := TCPHeader{
    Source:      0xaa47, // Random ephemeral port
    Destination: port,
    SeqNum:      rand.Uint32(),
    AckNum:      0,
    DataOffset:  5,      // 4 bits
    Reserved:    0,      // 3 bits
    ECN:         0,      // 3 bits
    Ctrl:        2,      // 6 bits (000010, SYN bit set)
    Window:      0xaaaa, // The amount of data that it is able to accept in bytes
    Checksum:    0,      // Kernel will set this if it's 0
    Urgent:      0,
    Options:     []TCPOption{},
  }

  data := packet.Marshal()
  packet.Checksum = Csum(data, to4byte(laddr), to4byte(raddr))

  data = packet.Marshal()

  conn, err := net.Dial("ip4:tcp", raddr)
  if err != nil {
    mlog.GetLogger().Fatal("Dial: %s\n", err)
  }

  sendTime := time.Now()

  numWrote, err := conn.Write(data)
  if err != nil {
    mlog.GetLogger().Fatalf("Write: %s\n", err)
  }
  if numWrote != len(data) {
    mlog.GetLogger().Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
  }

  conn.Close()

  return sendTime
}

func receiveTCPSynAck(localAddress, remoteAddress string) time.Time {
  netaddr, err := net.ResolveIPAddr("ip4", localAddress)
  if err != nil {
    mlog.GetLogger().Fatalf("net.ResolveIPAddr: %s. %s\n", localAddress, netaddr)
  }

  conn, err := net.ListenIP("ip4:tcp", netaddr)
  if err != nil {
    mlog.GetLogger().Fatalf("ListenIP: %s\n", err)
  }
  var receiveTime time.Time
  for {
    buf := make([]byte, 1024)
    numRead, raddr, err := conn.ReadFrom(buf)
    if err != nil {
      mlog.GetLogger().Fatalf("ReadFrom: %s\n", err)
    }
    if raddr.String() != remoteAddress {
      // this is not the packet we are looking for
      continue
    }
    receiveTime = time.Now()
    //fmt.Printf("Received: % x\n", buf[:numRead])
    tcp := NewTCPHeader(buf[:numRead])
    // Closed port gets RST, open port gets SYN ACK
    if tcp.HasFlag(RST) || (tcp.HasFlag(SYN) && tcp.HasFlag(ACK)) {
      break
    }
  }
  return receiveTime
}

// Parse packet into TCPHeader structure
func NewTCPHeader(data []byte) *TCPHeader {
  var tcp TCPHeader
  r := bytes.NewReader(data)
  binary.Read(r, binary.BigEndian, &tcp.Source)
  binary.Read(r, binary.BigEndian, &tcp.Destination)
  binary.Read(r, binary.BigEndian, &tcp.SeqNum)
  binary.Read(r, binary.BigEndian, &tcp.AckNum)

  var mix uint16
  binary.Read(r, binary.BigEndian, &mix)
  tcp.DataOffset = byte(mix >> 12)  // top 4 bits
  tcp.Reserved = byte(mix >> 9 & 7) // 3 bits
  tcp.ECN = byte(mix >> 6 & 7)      // 3 bits
  tcp.Ctrl = byte(mix & 0x3f)       // bottom 6 bits

  binary.Read(r, binary.BigEndian, &tcp.Window)
  binary.Read(r, binary.BigEndian, &tcp.Checksum)
  binary.Read(r, binary.BigEndian, &tcp.Urgent)

  return &tcp
}

func (tcp *TCPHeader) HasFlag(flagBit byte) bool {
  return tcp.Ctrl&flagBit != 0
}

func (tcp *TCPHeader) Marshal() []byte {

  buf := new(bytes.Buffer)
  binary.Write(buf, binary.BigEndian, tcp.Source)
  binary.Write(buf, binary.BigEndian, tcp.Destination)
  binary.Write(buf, binary.BigEndian, tcp.SeqNum)
  binary.Write(buf, binary.BigEndian, tcp.AckNum)

  var mix uint16
  mix = uint16(tcp.DataOffset)<<12 | // top 4 bits
    uint16(tcp.Reserved)<<9 | // 3 bits
    uint16(tcp.ECN)<<6 | // 3 bits
    uint16(tcp.Ctrl) // bottom 6 bits
  binary.Write(buf, binary.BigEndian, mix)

  binary.Write(buf, binary.BigEndian, tcp.Window)
  binary.Write(buf, binary.BigEndian, tcp.Checksum)
  binary.Write(buf, binary.BigEndian, tcp.Urgent)

  for _, option := range tcp.Options {
    binary.Write(buf, binary.BigEndian, option.Kind)
    if option.Length > 1 {
      binary.Write(buf, binary.BigEndian, option.Length)
      binary.Write(buf, binary.BigEndian, option.Data)
    }
  }

  out := buf.Bytes()

  // Pad to min tcp header size, which is 20 bytes (5 32-bit words)
  pad := 20 - len(out)
  for i := 0; i < pad; i++ {
    out = append(out, 0)
  }

  return out
}

// TCP Checksum
func Csum(data []byte, srcip, dstip [4]byte) uint16 {

  pseudoHeader := []byte{
    srcip[0], srcip[1], srcip[2], srcip[3],
    dstip[0], dstip[1], dstip[2], dstip[3],
    0,                  // zero
    6,                  // protocol number (6 == TCP)
    0, byte(len(data)), // TCP length (16 bits), not inc pseudo header
  }

  sumThis := make([]byte, 0, len(pseudoHeader)+len(data))
  sumThis = append(sumThis, pseudoHeader...)
  sumThis = append(sumThis, data...)
  //fmt.Printf("% x\n", sumThis)

  lenSumThis := len(sumThis)
  var nextWord uint16
  var sum uint32
  for i := 0; i+1 < lenSumThis; i += 2 {
    nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
    sum += uint32(nextWord)
  }
  if lenSumThis%2 != 0 {
    //fmt.Println("Odd byte")
    sum += uint32(sumThis[len(sumThis)-1])
  }

  // Add back any carry, and any carry from adding the carry
  sum = (sum >> 16) + (sum & 0xffff)
  sum = sum + (sum >> 16)

  // Bitwise complement
  return uint16(^sum)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
  c := make(chan struct{})
  go func() {
    defer close(c)
    wg.Wait()
  }()
  select {
  case <-c:
    return false // completed normally
  case <-time.After(timeout):
    return true // timed out
  }
}

func chooseInterface() string {
  interfaces, err := net.Interfaces()
  if err != nil {
    mlog.GetLogger().Fatalf("net.Interfaces: %s", err)
  }
  for _, iface := range interfaces {
    // Skip loopback
    if iface.Name == "lo" {
      continue
    }
    addrs, err := iface.Addrs()
    // Skip if error getting addresses
    if err != nil {
      mlog.GetLogger().Errorf("Error get addresses for interfaces %s. %s", iface.Name, err)
      continue
    }

    if len(addrs) > 0 {
      // This one will do
      return iface.Name
    }
  }
  return ""
}

func interfaceAddress(ifaceName string) net.Addr {
  iface, err := net.InterfaceByName(ifaceName)
  if err != nil {
    mlog.GetLogger().Fatalf("net.InterfaceByName for %s. %s", ifaceName, err)
  }
  addrs, err := iface.Addrs()
  if err != nil {
    mlog.GetLogger().Fatalf("iface.Addrs: %s", err)
  }
  return addrs[0]
}

func sendSyn(laddr, raddr string, port uint16) time.Time {

  packet := TCPHeader{
    Source:      0xaa47, // Random ephemeral port
    Destination: port,
    SeqNum:      rand.Uint32(),
    AckNum:      0,
    DataOffset:  5,      // 4 bits
    Reserved:    0,      // 3 bits
    ECN:         0,      // 3 bits
    Ctrl:        2,      // 6 bits (000010, SYN bit set)
    Window:      0xaaaa, // The amount of data that it is able to accept in bytes
    Checksum:    0,      // Kernel will set this if it's 0
    Urgent:      0,
    Options:     []TCPOption{},
  }

  data := packet.Marshal()
  packet.Checksum = Csum(data, to4byte(laddr), to4byte(raddr))

  data = packet.Marshal()

  fmt.Println("sendSyn: data: ", data)

  conn, err := net.Dial("ip4:tcp", raddr)
  if err != nil {
    mlog.GetLogger().Fatal("Dial: %s\n", err)
  }

  sendTime := time.Now()

  fmt.Println("sendSyn: sendTime: ", data)

  numWrote, err := conn.Write(data)
  if err != nil {
    mlog.GetLogger().Fatalf("Write: %s\n", err)
  }
  if numWrote != len(data) {
    mlog.GetLogger().Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
  }

  conn.Close()

  return sendTime
}

func to4byte(addr string) [4]byte {
  parts := strings.Split(addr, ".")
  b0, err := strconv.Atoi(parts[0])
  if err != nil {
    mlog.GetLogger().Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
  }
  b1, _ := strconv.Atoi(parts[1])
  b2, _ := strconv.Atoi(parts[2])
  b3, _ := strconv.Atoi(parts[3])
  return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}
