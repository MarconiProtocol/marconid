package mnet_ip

import (
  "reflect"
  "testing"
)

func TestGetMainInterfaceIpAddress(t *testing.T) {
  tests := []struct {
    name    string
    want    string
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := GetMainInterfaceIpAddress()
      if (err != nil) != tt.wantErr {
        t.Errorf("GetMainInterfaceIpAddress() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("GetMainInterfaceIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestGetOwnGatewayIpAddress(t *testing.T) {
  tests := []struct {
    name    string
    want    string
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := GetOwnGatewayIpAddress()
      if (err != nil) != tt.wantErr {
        t.Errorf("GetOwnGatewayIpAddress() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("GetOwnGatewayIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestGetCommandSetForAddIpAddressTunInterface(t *testing.T) {
  type args struct {
    taptunID      string
    ipAddr        string
    netmask       string
    peerIpAddr    string
    gatewayIpAddr string
  }
  tests := []struct {
    name  string
    args  args
    want  string
    want1 map[int][]string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, got1 := GetCommandSetForAddIpAddressTunInterface(tt.args.taptunID, tt.args.ipAddr, tt.args.netmask, tt.args.peerIpAddr, tt.args.gatewayIpAddr)
      if got != tt.want {
        t.Errorf("GetCommandSetForAddIpAddressTunInterface() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetCommandSetForAddIpAddressTunInterface() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}
