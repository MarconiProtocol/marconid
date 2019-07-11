package mnet_core_manager

import (
  "errors"
  "sync"
)

const (
  TUN_ID_BASE            uint16 = 10
  TUN_ID_MAX             uint16 = 1024
  SERVICE_BRIDGE_ID_BASE uint8  = 42
  SERVICE_BRIDGE_ID_MAX  uint8  = 60
)

/*
	Simple struct to keep track of id allocations for virtual network interfaces
*/
type InterfaceIDAllocations struct {
  serviceBridgeIds []bool
  connectionIds    []bool // true when allocated
  sync.Mutex
}

func initializeInterfaceIdAllocations() *InterfaceIDAllocations {
  ifIdAllocations := InterfaceIDAllocations{
    connectionIds:    make([]bool, TUN_ID_MAX-TUN_ID_BASE+1),
    serviceBridgeIds: make([]bool, SERVICE_BRIDGE_ID_MAX-SERVICE_BRIDGE_ID_BASE+1),
  }
  return &ifIdAllocations
}

/*
	Allocate the next available connection id and return it, otherwise returns an error
*/
func (nm *NetCoreManager) allocateNextConnectionId() (uint16, error) {
  nm.ifaceIDAllocations.Lock()
  defer nm.ifaceIDAllocations.Unlock()

  for i, allocated := range nm.ifaceIDAllocations.connectionIds {
    if !allocated {
      nm.ifaceIDAllocations.connectionIds[i] = true
      return uint16(i) + TUN_ID_BASE, nil
    }
  }
  return 0, errors.New("allocateNextConnectionId -> couldnt find an unallocated tun number")
}

/*
	Allocate the next available bridge interface id and return it, otherwise returns an error
*/
func (nm *NetCoreManager) allocateNextServiceBridgeId() (uint8, error) {
  nm.ifaceIDAllocations.Lock()
  defer nm.ifaceIDAllocations.Unlock()

  for i, allocated := range nm.ifaceIDAllocations.serviceBridgeIds {
    if !allocated {
      nm.ifaceIDAllocations.serviceBridgeIds[i] = true
      return uint8(i) + SERVICE_BRIDGE_ID_BASE, nil
    }
  }
  return 0, errors.New("allocateNextServiceBridgeId -> couldnt find an unallocated service bridge ID")
}

/*
	Deallocate the specified tun number
*/
func (nm *NetCoreManager) deallocateConnectionId(connectionId uint16) {
  nm.ifaceIDAllocations.Lock()
  defer nm.ifaceIDAllocations.Unlock()

  idx := connectionId - TUN_ID_BASE
  nm.ifaceIDAllocations.connectionIds[idx] = false
}
