package mruntime

import (
  "../../core"
  "runtime"
  "sync"
)

var instance *RuntimeInfo
var once sync.Once

type RuntimeInfo struct {
  osType        int
  InterfaceInfo *InterfaceRuntime
}

func GetMRuntime() *RuntimeInfo {
  once.Do(func() {
    instance = &RuntimeInfo{}
    instance.initialize()
  })
  return instance
}

func (r *RuntimeInfo) initialize() {
  r.InterfaceInfo = NewInterfaceRuntime()
}

func (r *RuntimeInfo) SetRuntimeOS() {
  osType, present := mcore.OSStringToInt[runtime.GOOS]
  if !present {
    osType = mcore.TYPE_OS_UNKNOWN
  }
  r.osType = osType
}

func (r *RuntimeInfo) GetRuntimeOS() int {
  return r.osType
}
