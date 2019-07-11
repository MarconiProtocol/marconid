package mnet_core_manager

import (
  "../../../peer"
  "../../vars"
  "reflect"
  "sync"
  "testing"
)

func TestGetNetCoreManager(t *testing.T) {
  tests := []struct {
    name string
    want *NetCoreManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetNetCoreManager(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetNetCoreManager() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_createNetCoreManager(t *testing.T) {
  tests := []struct {
    name string
    want *NetCoreManager
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := createNetCoreManager(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("createNetCoreManager() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestNetCoreManager_CreateMPipe(t *testing.T) {
  type fields struct {
    netIDtoBridgeInfoMap map[string]*BridgeInfo
    peerConnections      map[string]*Connection
    ifaceIDAllocations   *InterfaceIDAllocations
    pipeStates           *PipeStateMap
    tapCreationMutex     sync.Mutex
    tunCreationMutex     sync.Mutex
    Mutex                sync.Mutex
  }
  type args struct {
    args *mnet_vars.ConnectionArgs
  }
  tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      nm := &NetCoreManager{
        netIDtoBridgeInfoMap: tt.fields.netIDtoBridgeInfoMap,
        peerConnections:      tt.fields.peerConnections,
        ifaceIDAllocations:   tt.fields.ifaceIDAllocations,
        pipeStates:           tt.fields.pipeStates,
        tapCreationMutex:     tt.fields.tapCreationMutex,
        tunCreationMutex:     tt.fields.tunCreationMutex,
        Mutex:                tt.fields.Mutex,
      }
      if err := nm.CreateMPipe(tt.args.args); (err != nil) != tt.wantErr {
        t.Errorf("NetCoreManager.CreateMPipe() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestNetCoreManager_CloseMPipe(t *testing.T) {
  type fields struct {
    netIDtoBridgeInfoMap map[string]*BridgeInfo
    peerConnections      map[string]*Connection
    ifaceIDAllocations   *InterfaceIDAllocations
    pipeStates           *PipeStateMap
    tapCreationMutex     sync.Mutex
    tunCreationMutex     sync.Mutex
    Mutex                sync.Mutex
  }
  type args struct {
    peer *mpeer.Peer
  }
  tests := []struct {
    name   string
    fields fields
    args   args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      nm := &NetCoreManager{
        netIDtoBridgeInfoMap: tt.fields.netIDtoBridgeInfoMap,
        peerConnections:      tt.fields.peerConnections,
        ifaceIDAllocations:   tt.fields.ifaceIDAllocations,
        pipeStates:           tt.fields.pipeStates,
        tapCreationMutex:     tt.fields.tapCreationMutex,
        tunCreationMutex:     tt.fields.tunCreationMutex,
        Mutex:                tt.fields.Mutex,
      }
      nm.CloseMPipe(tt.args.peer)
    })
  }
}
