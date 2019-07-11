package mpacket_filter

import (
  "reflect"
  "sync"
  "testing"

  m_packet_filter "git.marconi.org/marconiprotocol/sdk/packet/filter"
  "github.com/google/gopacket"
)

func TestGetFilterManagerInstance(t *testing.T) {
  tests := []struct {
    name string
    want *PacketFilterManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetFilterManagerInstance(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetFilterManagerInstance() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestPacketFilterManager_init(t *testing.T) {
  tests := []struct {
    name          string
    filterManager *PacketFilterManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.filterManager.init()
    })
  }
}

func TestPacketFilterManager_loadConfig(t *testing.T) {
  tests := []struct {
    name          string
    filterManager *PacketFilterManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.filterManager.loadConfig()
    })
  }
}

func TestPacketFilterManager_loadFilter(t *testing.T) {
  tests := []struct {
    name          string
    filterManager *PacketFilterManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.filterManager.loadFilter()
    })
  }
}

func TestPacketFilterManager_ProcessPacket(t *testing.T) {
  type args struct {
    packet        *gopacket.Packet
    resultChannel *chan PacketFilterResult
  }
  tests := []struct {
    name          string
    filterManager *PacketFilterManager
    args          args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.filterManager.ProcessPacket(tt.args.packet, tt.args.resultChannel)
    })
  }
}

func Test_calculateFinalResult(t *testing.T) {
  type args struct {
    responses []m_packet_filter.FilterResponse
  }
  tests := []struct {
    name string
    args args
    want PacketFilterResult
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := calculateFinalResult(tt.args.responses...); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("calculateFinalResult() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_runFilter(t *testing.T) {
  type args struct {
    packet   gopacket.Packet
    filter   *m_packet_filter.Filter
    response *m_packet_filter.FilterResponse
    wg       *sync.WaitGroup
  }
  tests := []struct {
    name string
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runFilter(tt.args.packet, tt.args.filter, tt.args.response, tt.args.wg)
    })
  }
}
