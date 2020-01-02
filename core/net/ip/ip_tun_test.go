package mnet_ip

import (
  "reflect"
  "testing"
)

func TestConfigIpAddressForNewTunByCommand(t *testing.T) {
  type args struct {
    taptunID          string
    ipAddr            string
    netmask           string
    peerIpAddr        string
    peerGatewayIpAddr string
  }
  tests := []struct {
    name    string
    args    args
    wantRet bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if gotRet := ConfigIpAddressForNewTunByCommand(tt.args.taptunID, tt.args.ipAddr, tt.args.netmask, tt.args.peerIpAddr, tt.args.peerGatewayIpAddr); gotRet != tt.wantRet {
        t.Errorf("ConfigIpAddressForNewTunByCommand() = %v, want %v", gotRet, tt.wantRet)
      }
    })
  }
}

func TestConfigTunIpAddrByCommand(t *testing.T) {
  type args struct {
    taptunID          string
    ipAddr            string
    netmask           string
    peerIpAddr        string
    peerGatewayIpAddr string
  }
  tests := []struct {
    name       string
    args       args
    wantResult map[int]string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if gotResult := ConfigTunIpAddrByCommand(tt.args.taptunID, tt.args.ipAddr, tt.args.netmask, tt.args.peerIpAddr, tt.args.peerGatewayIpAddr); !reflect.DeepEqual(gotResult, tt.wantResult) {
        t.Errorf("ConfigTunIpAddrByCommand() = %v, want %v", gotResult, tt.wantResult)
      }
    })
  }
}

func TestConfigMconnIpAddress(t *testing.T) {
  type args struct {
    taptunNum  string
    ipAddr     string
    netmask    string
    gwIpAddr   string
    peerIpAddr string
  }
  tests := []struct {
    name string
    args args
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      ConfigMconnIpAddress(tt.args.taptunNum, tt.args.ipAddr, tt.args.netmask, tt.args.gwIpAddr, tt.args.peerIpAddr)
    })
  }
}

func Test_configClientLayerIpAddress(t *testing.T) {
  type args struct {
    taptunID          string
    ipAddr            string
    netmask           string
    peerIpAddr        string
    peerGatewayIpAddr string
  }
  tests := []struct {
    name string
    args args
    want bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := configClientLayerIpAddress(tt.args.taptunID, tt.args.ipAddr, tt.args.netmask, tt.args.peerIpAddr, tt.args.peerGatewayIpAddr); got != tt.want {
        t.Errorf("configClientLayerIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_allowTrafficForwardOnTunInterface(t *testing.T) {
  type args struct {
    taptunID string
  }
  tests := []struct {
    name       string
    args       args
    wantResult map[int]string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if gotResult := allowTrafficForwardOnTunInterface(tt.args.taptunID); !reflect.DeepEqual(gotResult, tt.wantResult) {
        t.Errorf("allowTrafficForwardOnTunInterface() = %v, want %v", gotResult, tt.wantResult)
      }
    })
  }
}

func Test_allowTrafficForwardOnSystem(t *testing.T) {
  tests := []struct {
    name       string
    wantResult map[int]string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if gotResult := allowTrafficForwardOnSystem(); !reflect.DeepEqual(gotResult, tt.wantResult) {
        t.Errorf("allowTrafficForwardOnSystem() = %v, want %v", gotResult, tt.wantResult)
      }
    })
  }
}
