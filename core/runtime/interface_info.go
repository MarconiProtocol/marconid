package mruntime

import (
  "sync"
)

//Tap or Tun software interface runtime info
type InterfaceRuntime struct {
  localMainInterfaceIpAddr string
  sync.Mutex
}

func NewInterfaceRuntime() *InterfaceRuntime {
  ifr := &InterfaceRuntime{}
  return ifr
}

// localMainInterfaceIpAddr
func (i *InterfaceRuntime) GetLocalMainInterfaceIpAddr() string {
  return i.localMainInterfaceIpAddr
}
func (i *InterfaceRuntime) SetLocalMainInterfaceIpAddr(newLocalMainInterfaceIpAddr string) {
  i.localMainInterfaceIpAddr = newLocalMainInterfaceIpAddr
}
