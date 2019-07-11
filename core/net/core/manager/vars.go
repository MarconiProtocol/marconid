package mnet_core_manager

import "sync"

type PipeState int

const (
  UNATTEMPTED PipeState = iota // 0
  ATTEMPTING                   // 1
  SUCCESS                      // 2
)

type PipeStateMap struct {
  StateMap *map[string]PipeState
  sync.Mutex
}

func NewPipeStates() *PipeStateMap {
  psm := PipeStateMap{}
  states := make(map[string]PipeState)
  psm.StateMap = &states
  return &psm
}

type BridgeInfo struct {
  ID string
  IP string
}
