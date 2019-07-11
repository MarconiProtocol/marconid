package mnet_if

import (
  "reflect"
  "testing"
)

func TestCreateInterface(t *testing.T) {
  type args struct {
    fd     uintptr
    ifName string
    flags  uint16
  }
  tests := []struct {
    name              string
    args              args
    wantCreatedIFName string
    wantErr           bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotCreatedIFName, err := CreateInterface(tt.args.fd, tt.args.ifName, tt.args.flags)
      if (err != nil) != tt.wantErr {
        t.Errorf("CreateInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotCreatedIFName != tt.wantCreatedIFName {
        t.Errorf("CreateInterface() = %v, want %v", gotCreatedIFName, tt.wantCreatedIFName)
      }
    })
  }
}

func Test_newTUN(t *testing.T) {
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
      gotIfce, err := newTUN(tt.args.ifName)
      if (err != nil) != tt.wantErr {
        t.Errorf("newTUN() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(gotIfce, tt.wantIfce) {
        t.Errorf("newTUN() = %v, want %v", gotIfce, tt.wantIfce)
      }
    })
  }
}

func Test_newTAP(t *testing.T) {
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
      gotIfce, err := newTAP(tt.args.ifName)
      if (err != nil) != tt.wantErr {
        t.Errorf("newTAP() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(gotIfce, tt.wantIfce) {
        t.Errorf("newTAP() = %v, want %v", gotIfce, tt.wantIfce)
      }
    })
  }
}

func Test_createInterface(t *testing.T) {
  type args struct {
    fd     uintptr
    ifName string
    flags  uint16
  }
  tests := []struct {
    name              string
    args              args
    wantCreatedIFName string
    wantErr           bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotCreatedIFName, err := createInterface(tt.args.fd, tt.args.ifName, tt.args.flags)
      if (err != nil) != tt.wantErr {
        t.Errorf("createInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotCreatedIFName != tt.wantCreatedIFName {
        t.Errorf("createInterface() = %v, want %v", gotCreatedIFName, tt.wantCreatedIFName)
      }
    })
  }
}

func TestInterface_Open2(t *testing.T) {
  type args struct {
    mtu    uint
    ifNum  string
    ifType string
  }
  tests := []struct {
    name    string
    vni     *Interface
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.vni.Open2(tt.args.mtu, tt.args.ifNum, tt.args.ifType); (err != nil) != tt.wantErr {
        t.Errorf("Interface.Open2() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestInterface_Open(t *testing.T) {
  type args struct {
    mtu    uint
    tapNum string
    ifType string
  }
  tests := []struct {
    name     string
    tap_conn *Interface
    args     args
    wantErr  bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := tt.tap_conn.Open(tt.args.mtu, tt.args.tapNum, tt.args.ifType); (err != nil) != tt.wantErr {
        t.Errorf("Interface.Open() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestInterface_Close(t *testing.T) {
  tests := []struct {
    name string
    inet *Interface
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.inet.Close()
    })
  }
}

func TestInterface_Read(t *testing.T) {
  type args struct {
    b []byte
  }
  tests := []struct {
    name    string
    inet    *Interface
    args    args
    wantN   int
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotN, err := tt.inet.Read(tt.args.b)
      if (err != nil) != tt.wantErr {
        t.Errorf("Interface.Read() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotN != tt.wantN {
        t.Errorf("Interface.Read() = %v, want %v", gotN, tt.wantN)
      }
    })
  }
}

func TestInterface_Write(t *testing.T) {
  type args struct {
    b []byte
  }
  tests := []struct {
    name    string
    inet    *Interface
    args    args
    wantN   int
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotN, err := tt.inet.Write(tt.args.b)
      if (err != nil) != tt.wantErr {
        t.Errorf("Interface.Write() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotN != tt.wantN {
        t.Errorf("Interface.Write() = %v, want %v", gotN, tt.wantN)
      }
    })
  }
}
