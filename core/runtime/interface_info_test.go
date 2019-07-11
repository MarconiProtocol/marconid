package mruntime

import (
  "reflect"
  "testing"
)

func TestNewInterfaceRuntime(t *testing.T) {
  tests := []struct {
    name string
    want *InterfaceRuntime
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := NewInterfaceRuntime(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("NewInterfaceRuntime() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterfaceRuntime_GetLocalMainInterfaceIpAddr(t *testing.T) {
  tests := []struct {
    name string
    i    *InterfaceRuntime
    want string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.i.GetLocalMainInterfaceIpAddr(); got != tt.want {
        t.Errorf("InterfaceRuntime.GetLocalMainInterfaceIpAddr() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterfaceRuntime_SetLocalMainInterfaceIpAddr(t *testing.T) {
  type args struct {
    newLocalMainInterfaceIpAddr string
  }
  tests := []struct {
    name string
    i    *InterfaceRuntime
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.i.SetLocalMainInterfaceIpAddr(tt.args.newLocalMainInterfaceIpAddr)
    })
  }
}
