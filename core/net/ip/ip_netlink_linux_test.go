package mnet_ip

import (
  "net"
  "reflect"
  "testing"

  "gitlab.neji.vm.tc/marconi/netlink"
)

func Test_getNetlinkAddress(t *testing.T) {
  type args struct {
    ipAddr  string
    netmask int
  }
  tests := []struct {
    name string
    args args
    want *netlink.Addr
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := getNetlinkAddress(tt.args.ipAddr, tt.args.netmask); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("getNetlinkAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_getNetlinkLink(t *testing.T) {
  type args struct {
    ifName string
  }
  tests := []struct {
    name string
    args args
    want netlink.Link
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := getNetlinkLink(tt.args.ifName); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("getNetlinkLink() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_getAllLinksByNetlink(t *testing.T) {
  tests := []struct {
    name string
    want []netlink.Link
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := getAllLinksByNetlink(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("getAllLinksByNetlink() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_getIPNetAddress(t *testing.T) {
  type args struct {
    ipAddr  string
    netmask int
  }
  tests := []struct {
    name string
    args args
    want *net.IPNet
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := getIPNetAddress(tt.args.ipAddr, tt.args.netmask); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("getIPNetAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_getNetIPv4(t *testing.T) {
  type args struct {
    ipAddr string
  }
  tests := []struct {
    name string
    args args
    want net.IP
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := getNetIPv4(tt.args.ipAddr); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("getNetIPv4() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestAssignNetlinkIpAddress(t *testing.T) {
  type args struct {
    ifName  string
    ipAddr  string
    netmask int
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      AssignNetlinkIpAddress(tt.args.ifName, tt.args.ipAddr, tt.args.netmask)
    })
  }
}

func TestRemoveNetlinkIpAddress(t *testing.T) {
  type args struct {
    ifName  string
    ipAddr  string
    netmask int
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      RemoveNetlinkIpAddress(tt.args.ifName, tt.args.ipAddr, tt.args.netmask)
    })
  }
}

func TestGetNetlinkAllIpAddress(t *testing.T) {
  type args struct {
    ifName string
    family int
  }
  tests := []struct {
    name string
    args args
    want []netlink.Addr
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetNetlinkAllIpAddress(tt.args.ifName, tt.args.family); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetNetlinkAllIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestGetIpAddressV4ByNetlink(t *testing.T) {
  type args struct {
    ifName string
  }
  tests := []struct {
    name string
    args args
    want []netlink.Addr
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetIpAddressV4ByNetlink(tt.args.ifName); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetIpAddressV4ByNetlink() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestGetIpAddressV6ByNetlink(t *testing.T) {
  type args struct {
    ifName string
  }
  tests := []struct {
    name string
    args args
    want []netlink.Addr
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetIpAddressV6ByNetlink(tt.args.ifName); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetIpAddressV6ByNetlink() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestListNetlinkAllIpAddress(t *testing.T) {
  type args struct {
    ifName string
    family int
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      ListNetlinkAllIpAddress(tt.args.ifName, tt.args.family)
    })
  }
}

func TestListAllLinksByNetlink(t *testing.T) {
  tests := []struct {
    name string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      ListAllLinksByNetlink()
    })
  }
}
