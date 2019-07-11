package mnet_core_udp

import (
  "../../../if"
  "net"
  "reflect"
  "testing"
)

func TestGetUDPTransport(t *testing.T) {
  tests := []struct {
    name string
    want *UDPTransport
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetUDPTransport(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetUDPTransport() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestUDPTransport_ListenAndTransmit(t *testing.T) {
  type args struct {
    localIpAddr           string
    localPort             string
    remoteIpAddr          string
    remotePort            string
    tapConn               *mnet_if.Interface
    key                   []byte
    dataKey               *[]byte
    isSecure              bool
    isTun                 bool
    listenSignalChannel   chan string
    transmitSignalChannel chan string
  }
  tests := []struct {
    name    string
    udpt    *UDPTransport
    args    args
    want    net.Conn
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      udpt := &UDPTransport{}
      got, err := udpt.ListenAndTransmit(tt.args.localIpAddr, tt.args.localPort, tt.args.remoteIpAddr, tt.args.remotePort, tt.args.tapConn, tt.args.key, tt.args.dataKey, tt.args.isSecure, tt.args.isTun, tt.args.listenSignalChannel, tt.args.transmitSignalChannel)
      if (err != nil) != tt.wantErr {
        t.Errorf("UDPTransport.ListenAndTransmit() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UDPTransport.ListenAndTransmit() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestUDPTransport_Listen(t *testing.T) {
  type args struct {
    physConn             *net.UDPConn
    tapConn              *mnet_if.Interface
    peerAddr             *net.UDPAddr
    key                  []byte
    dataKey              *[]byte
    isSecure             bool
    peerDiscoveryChannel chan net.UDPAddr
    isTun                bool
    signalChannel        chan string
  }
  tests := []struct {
    name string
    udpt *UDPTransport
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      udpt := &UDPTransport{}
      udpt.Listen(tt.args.physConn, tt.args.tapConn, tt.args.peerAddr, tt.args.key, tt.args.dataKey, tt.args.isSecure, tt.args.peerDiscoveryChannel, tt.args.isTun, tt.args.signalChannel)
    })
  }
}

func TestUDPTransport_Transmit(t *testing.T) {
  type args struct {
    physConn             *net.UDPConn
    tapConn              *mnet_if.Interface
    peerAddr             *net.UDPAddr
    key                  []byte
    dataKey              *[]byte
    isSecure             bool
    peerDiscoveryChannel chan net.UDPAddr
    isTun                bool
    signalChannel        chan string
  }
  tests := []struct {
    name string
    udpt *UDPTransport
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      udpt := &UDPTransport{}
      udpt.Transmit(tt.args.physConn, tt.args.tapConn, tt.args.peerAddr, tt.args.key, tt.args.dataKey, tt.args.isSecure, tt.args.peerDiscoveryChannel, tt.args.isTun, tt.args.signalChannel)
    })
  }
}
