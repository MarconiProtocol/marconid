package mruntime

import (
  "reflect"
  "testing"
)

func TestGetMRuntime(t *testing.T) {
  tests := []struct {
    name string
    want *RuntimeInfo
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetMRuntime(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetMRuntime() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestRuntimeInfo_initialize(t *testing.T) {
  tests := []struct {
    name string
    r    *RuntimeInfo
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.r.initialize()
    })
  }
}

func TestRuntimeInfo_SetRuntimeOS(t *testing.T) {
  tests := []struct {
    name string
    r    *RuntimeInfo
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.r.SetRuntimeOS()
    })
  }
}

func TestRuntimeInfo_GetRuntimeOS(t *testing.T) {
  tests := []struct {
    name string
    r    *RuntimeInfo
    want int
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.r.GetRuntimeOS(); got != tt.want {
        t.Errorf("RuntimeInfo.GetRuntimeOS() = %v, want %v", got, tt.want)
      }
    })
  }
}
