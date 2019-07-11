package mnet_ip

import "testing"

func TestConfigBridgeIpAddrByCommand(t *testing.T) {
  type args struct {
    bridgeID string
    ipAddr   string
    netmask  string
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := ConfigBridgeIpAddrByCommand(tt.args.bridgeID, tt.args.ipAddr, tt.args.netmask)
      if (err != nil) != tt.wantErr {
        t.Errorf("ConfigBridgeIpAddrByCommand() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ConfigBridgeIpAddrByCommand() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestConfigBridgeUpByCommand(t *testing.T) {
  type args struct {
    bridgeID string
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := ConfigBridgeUpByCommand(tt.args.bridgeID)
      if (err != nil) != tt.wantErr {
        t.Errorf("ConfigBridgeUpByCommand() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ConfigBridgeUpByCommand() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestConfigTapToBridgeByCommand(t *testing.T) {
  type args struct {
    bridgeID string
    tapID    string
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := ConfigTapToBridgeByCommand(tt.args.bridgeID, tt.args.tapID)
      if (err != nil) != tt.wantErr {
        t.Errorf("ConfigTapToBridgeByCommand() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ConfigTapToBridgeByCommand() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestConfigBridgeByCommand(t *testing.T) {
  type args struct {
    bridgeID    string
    ipAddr      string
    netmask     string
    resetBridge bool
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := ConfigBridgeByCommand(tt.args.bridgeID, tt.args.ipAddr, tt.args.netmask, tt.args.resetBridge)
      if (err != nil) != tt.wantErr {
        t.Errorf("ConfigBridgeByCommand() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ConfigBridgeByCommand() = %v, want %v", got, tt.want)
      }
    })
  }
}
