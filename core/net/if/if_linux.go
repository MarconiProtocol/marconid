package mnet_if

import (
  "../vars"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "os"
  "strings"
  "syscall"
  "unsafe"
)

func CreateInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
  var req IfReq
  req.Flags = flags
  copy(req.Name[:], ifName)
  _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
  if errno != 0 {
    err = errno
    return
  }
  createdIFName = strings.Trim(string(req.Name[:]), "\x00")
  return
}

func newTUN(ifName string) (ifce *Interface, err error) {
  file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
  if err != nil {
    mlog.GetLogger().Error("newTUN os.OpenFile failed: ", err)
    return nil, err
  }
  name, err := CreateInterface(file.Fd(), ifName, cIFF_TUN|cIFF_NO_PI)
  if err != nil {
    mlog.GetLogger().Error("newTUN CreateInterface failed: ", err)
    return nil, err
  }
  ifce = &Interface{isTAP: false, ReadWriteCloser: file, name: name}
  return
}

func newTAP(ifName string) (ifce *Interface, err error) {
  file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
  if err != nil {
    mlog.GetLogger().Error("newTAP os.OpenFile failed: ", err)
    return nil, err
  }
  name, err := CreateInterface(file.Fd(), ifName, cIFF_TAP|cIFF_NO_PI)
  if err != nil {
    mlog.GetLogger().Error("newTAP CreateInterface failed: ", err)
    return nil, err
  }
  ifce = &Interface{isTAP: true, ReadWriteCloser: file, name: name}
  return
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
  var req IfReq
  req.Flags = flags
  copy(req.Name[:], ifName)
  _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
  if errno != 0 {
    err = errno
    return
  }
  createdIFName = strings.Trim(string(req.Name[:]), "\x00")
  return
}

//for tun
//new open for tun device
//currently code goes into this path for linux tun connection
//TODO: move it to regular Open()
func (vni *Interface) Open2(mtu uint, ifNum string, ifType string) (err error) {
  vni.fd_os, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
  ifName := ifType + ifNum
  if err != nil {
    mlog.GetLogger().Error("Open2 os.OpenFile failed: ", err)
    return
  }
  name, err := createInterface(vni.fd_os.Fd(), ifName, cIFF_TUN|cIFF_NO_PI)
  if err != nil {
    return
  }
  vni.ifname = name
  mlog.GetLogger().Info("Open2 ifType: ", ifType, ", name: ", name)
  return
}

/**********************************************************************/
/*** Tap Device Open/Close/Read/Write ***/
/**********************************************************************/
//NOTE: legacy tap device create method
//used with in tap driver only for now
func (tap_conn *Interface) Open(mtu uint, tapNum string, ifType string) (err error) {
  /* Open the tap/tun device */
  tap_conn.fd, err = syscall.Open("/dev/net/tun", syscall.O_RDWR, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IROTH)
  if err != nil {
    mlog.GetLogger().Error("Open failed to open device /dev/net/tun: ", err)
    return fmt.Errorf("Error opening device /dev/net/tun: %s", err)
  }

  //NOTE: below make l2/tap drivers
  /* Prepare a struct ifreq structure for TUNSETIFF with tap settings */
  /* IFF_TAP: tap device, IFF_NO_PI: no extra packet information */
  var ifr_flags uint16 = 0
  if ifType == "tap" {
    ifr_flags = uint16(syscall.IFF_TAP | syscall.IFF_NO_PI)
  } else if ifType == "tun" {
    ifr_flags = uint16(syscall.IFF_TUN | syscall.IFF_NO_PI)
  } else {
    ifr_flags = uint16(syscall.IFF_TAP | syscall.IFF_NO_PI)
  }

  /* FIXME: Assumes little endian */
  ifr_struct := make([]byte, 32)
  ifr_struct[16] = byte(ifr_flags)
  ifr_struct[17] = byte(ifr_flags >> 8)

  //NOTE: override 1st 10 bytes on ifr_struct
  //LocalMetadata is preset by user or system if need to set if name
  copy(ifr_struct[:], mnet_vars.MPIPE_LINK_DEVICE_NAME_PREFIX+tapNum)

  r0, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_conn.fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifr_struct[0])))
  if r0 != 0 {
    tap_conn.Close()
    return fmt.Errorf("Error setting tun/tap type: %s", err)
  }

  /* Extract the assigned interface name into a string */
  tap_conn.ifname = string(ifr_struct[0:16])

  /* Create a raw socket for our tap interface, so we can set the MTU */
  tap_sockfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
  if err != nil {
    mlog.GetLogger().Error("Creating packet socket failed: ", err)
    tap_conn.Close()
    return fmt.Errorf("Error creating packet socket: %s", err)
  }
  /* We won't need the socket after we've set the MTU and brought the
   * interface up */
  defer syscall.Close(tap_sockfd)

  /* Bind the raw socket to our tap interface */
  err = syscall.BindToDevice(tap_sockfd, tap_conn.ifname)
  if err != nil {
    mlog.GetLogger().Error("Binding packet socket to tap interface failed: ", err)
    tap_conn.Close()
    return fmt.Errorf("Error binding packet socket to tap interface: %s", err)
  }

  /* Prepare a ifreq structure for SIOCSIFMTU with MTU setting */
  ifr_mtu := mtu
  /* FIXME: Assumes little endian */
  ifr_struct[16] = byte(ifr_mtu)
  ifr_struct[17] = byte(ifr_mtu >> 8)
  ifr_struct[18] = byte(ifr_mtu >> 16)
  ifr_struct[19] = byte(ifr_mtu >> 24)
  r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCSIFMTU, uintptr(unsafe.Pointer(&ifr_struct[0])))
  if r0 != 0 {
    mlog.GetLogger().Error("Error setting MTU on tap interface: ", err)
    tap_conn.Close()
    return fmt.Errorf("Error setting MTU on tap interface: %s", err)
  }

  /* Get the current interface flags in ifr_struct */
  r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifr_struct[0])))
  if r0 != 0 {
    mlog.GetLogger().Error("Error getting tap interface flags: ", err)
    tap_conn.Close()
    return fmt.Errorf("Error getting tap interface flags: %s", err)
  }
  /* Update the interface flags to bring the interface up */
  /* FIXME: Assumes little endian */
  ifr_flags = uint16(ifr_struct[16]) | (uint16(ifr_struct[17]) << 8)
  ifr_flags |= syscall.IFF_UP | syscall.IFF_RUNNING
  ifr_struct[16] = byte(ifr_flags)
  ifr_struct[17] = byte(ifr_flags >> 8)
  r0, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap_sockfd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifr_struct[0])))
  if r0 != 0 {
    mlog.GetLogger().Error("Error bringing tap interface: ", err)
    tap_conn.Close()
    return fmt.Errorf("Error bringing up tap interface: %s", err)
  }
  return nil
}

func (inet *Interface) Close() {
  if inet.isTAP {
    syscall.Close(inet.fd)
  } else {
    //tun
    inet.fd_os.Close()
  }
}

func (inet *Interface) Read(b []byte) (n int, err error) {
  if inet.isTAP {
    return syscall.Read(inet.fd, b)
  } else {
    //tun
    return inet.fd_os.Read(b)
  }
}

func (inet *Interface) Write(b []byte) (n int, err error) {
  if inet.isTAP {
    return syscall.Write(inet.fd, b)
  } else {
    //tun
    return inet.fd_os.Write(b)
  }
}
