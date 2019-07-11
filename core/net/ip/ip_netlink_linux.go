package mnet_ip

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
  "gitlab.neji.vm.tc/marconi/netlink"
  "golang.org/x/sys/unix"
  "net"
  "strconv"
  "strings"
)

const (
  // Family type definitions
  FAMILY_ALL  = unix.AF_UNSPEC
  FAMILY_V4   = unix.AF_INET
  FAMILY_V6   = unix.AF_INET6
  FAMILY_MPLS = 28
  // Arbitrary set value (greater than default 4k) to allow receiving
  // from kernel more verbose messages e.g. for statistics,
  // tc rules or filters, or other more memory requiring data.
  RECEIVE_BUFFER_SIZE = 65536
)

//// netlink
//used with ip assignment to interfaces
func getNetlinkAddress(ipAddr string, netmask int) *netlink.Addr {
  address := getIPNetAddress(ipAddr, netmask)
  addr := &netlink.Addr{IPNet: address}
  return addr
}

func getNetlinkLink(ifName string) netlink.Link {
  link, err := netlink.LinkByName(ifName)
  if err != nil {
    mlog.GetLogger().Warn("getNetlinkLink failed", err)
  }
  return link
}

func getAllLinksByNetlink() []netlink.Link {
  netlink.LinkList()
  links, err := netlink.LinkList()
  if err != nil {
    mlog.GetLogger().Warn("getNetlinkLink failed", err)
  }
  return links
}

//// ip.net
func getIPNetAddress(ipAddr string, netmask int) *net.IPNet {
  //IPv4 only for now
  ip := getNetIPv4(ipAddr)
  address := &net.IPNet{IP: ip, Mask: net.CIDRMask(netmask, 32)}
  return address
}

func getNetIPv4(ipAddr string) net.IP {
  addressHex := strings.Split(ipAddr, ".")
  if len(addressHex) == 4 {
    a, _ := strconv.Atoi(addressHex[0])
    b, _ := strconv.Atoi(addressHex[1])
    c, _ := strconv.Atoi(addressHex[2])
    d, _ := strconv.Atoi(addressHex[3])
    ip := net.IPv4(byte(a), byte(b), byte(c), byte(d))
    return ip
  }
  return nil
}

func AssignNetlinkIpAddress(ifName string, ipAddr string, netmask int) {
  address := getNetlinkAddress(ipAddr, netmask)
  link := getNetlinkLink(ifName)

  err := netlink.AddrAdd(link, address)
  if err != nil {
    mlog.GetLogger().Error("Replacing IP Failed.", err)
  }
  mlog.GetLogger().Info(fmt.Sprintf("Assigned IP %s/%d to %s", ipAddr, netmask, ifName))
}

func RemoveNetlinkIpAddress(ifName string, ipAddr string, netmask int) {
  address := getNetlinkAddress(ipAddr, netmask)
  link := getNetlinkLink(ifName)
  err := netlink.AddrDel(link, address)
  if err != nil {
    mlog.GetLogger().Error("Deleting IP Failed.", err)
  }
}

func GetNetlinkAllIpAddress(ifName string, family int) []netlink.Addr {
  link := getNetlinkLink(ifName)
  addrs, err := netlink.AddrList(link, family)
  if err != nil {
    mlog.GetLogger().Error("GetNetlinkAllIpAddress failed", err)
  }
  return addrs
}

func GetIpAddressV4ByNetlink(ifName string) []netlink.Addr {
  return GetNetlinkAllIpAddress(ifName, netlink.FAMILY_V4)
}

func GetIpAddressV6ByNetlink(ifName string) []netlink.Addr {
  return GetNetlinkAllIpAddress(ifName, netlink.FAMILY_V6)
}

func ListNetlinkAllIpAddress(ifName string, family int) {
  fmt.Println(GetNetlinkAllIpAddress(ifName, family))
}

func ListAllLinksByNetlink() {
  links := getAllLinksByNetlink()
  for i, link := range links {
    fmt.Println(strconv.Itoa(i) + ": " + link.Attrs().Name)
  }
}
