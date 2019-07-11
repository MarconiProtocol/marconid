package mnet_if

import (
  "io"
  "math/rand"
  "os"
  "time"
)

const (
  cIFF_TUN   = 0x0001
  cIFF_TAP   = 0x0002
  cIFF_NO_PI = 0x1000

  IPv6_HEADER_LENGTH = 40
)

// Interface is a TUN/TAP interface.
type Interface struct {
  //TODO: this needs to be unified as fd either using syscall or os calls.
  fd_os  *os.File
  fd_sys uintptr
  fd     int
  ifname string
  isTAP  bool
  io.ReadWriteCloser
  name string
  mtu  uint
}

//// Tap Device Open/Close/Read/Write
type TapConn struct {
  fd     int
  ifname string
}

type IfReq struct {
  Name  [0x10]byte
  Flags uint16
  pad   [0x28 - 0x10 - 2]byte
}

func (netIf *Interface) OpenTun(mtu uint, tapNum string) (err error) {
  return netIf.OpenTunInterface(mtu, tapNum)
}

func (netIf *Interface) OpenTap(mtu uint, tapNum string) (err error) {
  return netIf.OpenTapInterface(mtu, tapNum)
}

//new
func (netIf *Interface) OpenTunInterface(mtu uint, tapNum string) (err error) {
  //TODO: it needs to be implemented
  //TODO: new
  return netIf.Open2(mtu, tapNum, "tun")
}

//new
func (netIf *Interface) OpenTapInterface(mtu uint, tapNum string) (err error) {
  return netIf.Open(mtu, tapNum, "tap")
}

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTAP(ifName string) (ifce *Interface, err error) {
  return newTAP(ifName)
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTUN(ifName string) (ifce *Interface, err error) {
  return newTUN(ifName)
}

// Returns true if ifce is a TUN interface, otherwise returns false;
func (ifce *Interface) IsTUN() bool {
  return !ifce.isTAP
}

// Returns true if ifce is a TAP interface, otherwise returns false;
func (ifce *Interface) IsTAP() bool {
  return ifce.isTAP
}

func (ifce *Interface) GetName() string {
  return ifce.ifname
}

func (ifce *Interface) GetFdOS() *os.File {
  return ifce.fd_os
}

func (ifce *Interface) GetFd() int {
  return ifce.fd
}

func (ifce *Interface) SetTap(isTap bool) {
  ifce.isTAP = isTap
}

// Returns the interface name of ifce, e.g. tun0, tap1, etc..
func (ifce *Interface) Name() string {
  return ifce.name
}

func TunIdRandom(min, max int) int {
  rand.Seed(time.Now().Unix())
  return rand.Intn(max-min) + min
}
