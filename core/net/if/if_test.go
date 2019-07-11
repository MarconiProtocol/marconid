package mnet_if

import (
  "os"
  "reflect"
  "testing"
)

func TestInterface_OpenTun(t *testing.T) {
  type args struct {
    mtu    uint
    tapNum string
  }
  tests := []struct {
    name    string
    netIf   *Interface
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.netIf.OpenTun(tt.args.mtu, tt.args.tapNum); (err != nil) != tt.wantErr {
        t.Errorf("Interface.OpenTun() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestInterface_OpenTap(t *testing.T) {
  type args struct {
    mtu    uint
    tapNum string
  }
  tests := []struct {
    name    string
    netIf   *Interface
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.netIf.OpenTap(tt.args.mtu, tt.args.tapNum); (err != nil) != tt.wantErr {
        t.Errorf("Interface.OpenTap() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestInterface_OpenTunInterface(t *testing.T) {
  type args struct {
    mtu    uint
    tapNum string
  }
  tests := []struct {
    name    string
    netIf   *Interface
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.netIf.OpenTunInterface(tt.args.mtu, tt.args.tapNum); (err != nil) != tt.wantErr {
        t.Errorf("Interface.OpenTunInterface() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestInterface_OpenTapInterface(t *testing.T) {
  type args struct {
    mtu    uint
    tapNum string
  }
  tests := []struct {
    name    string
    netIf   *Interface
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.netIf.OpenTapInterface(tt.args.mtu, tt.args.tapNum); (err != nil) != tt.wantErr {
        t.Errorf("Interface.OpenTapInterface() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestNewTAP(t *testing.T) {
  type args struct {
    ifName string
  }
  tests := []struct {
    name     string
    args     args
    wantIfce *Interface
    wantErr  bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotIfce, err := NewTAP(tt.args.ifName)
      if (err != nil) != tt.wantErr {
        t.Errorf("NewTAP() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(gotIfce, tt.wantIfce) {
        t.Errorf("NewTAP() = %v, want %v", gotIfce, tt.wantIfce)
      }
    })
  }
}

func TestNewTUN(t *testing.T) {
  type args struct {
    ifName string
  }
  tests := []struct {
    name     string
    args     args
    wantIfce *Interface
    wantErr  bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotIfce, err := NewTUN(tt.args.ifName)
      if (err != nil) != tt.wantErr {
        t.Errorf("NewTUN() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(gotIfce, tt.wantIfce) {
        t.Errorf("NewTUN() = %v, want %v", gotIfce, tt.wantIfce)
      }
    })
  }
}

func TestInterface_IsTUN(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.IsTUN(); got != tt.want {
        t.Errorf("Interface.IsTUN() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterface_IsTAP(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.IsTAP(); got != tt.want {
        t.Errorf("Interface.IsTAP() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterface_GetName(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.GetName(); got != tt.want {
        t.Errorf("Interface.GetName() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterface_GetFdOS(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want *os.File
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.GetFdOS(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Interface.GetFdOS() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterface_GetFd(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want int
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.GetFd(); got != tt.want {
        t.Errorf("Interface.GetFd() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestInterface_SetTap(t *testing.T) {
  type args struct {
    isTap bool
  }
  tests := []struct {
    name string
    ifce *Interface
    args args
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.ifce.SetTap(tt.args.isTap)
    })
  }
}

func TestInterface_Name(t *testing.T) {
  tests := []struct {
    name string
    ifce *Interface
    want string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.ifce.Name(); got != tt.want {
        t.Errorf("Interface.Name() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestTunIdRandom(t *testing.T) {
  type args struct {
    min int
    max int
  }
  tests := []struct {
    name string
    args args
    want int
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := TunIdRandom(tt.args.min, tt.args.max); got != tt.want {
        t.Errorf("TunIdRandom() = %v, want %v", got, tt.want)
      }
    })
  }
}
