package msys_cmd_centos7

import "testing"

func TestCmdSuite_GetOwnGatewayIpAddress(t *testing.T) {
  tests := []struct {
    name    string
    c       CmdSuite
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.GetOwnGatewayIpAddress()
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.GetOwnGatewayIpAddress() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.GetOwnGatewayIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_GetMainInterfaceIpAddress(t *testing.T) {
  tests := []struct {
    name    string
    c       CmdSuite
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.GetMainInterfaceIpAddress()
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.GetMainInterfaceIpAddress() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.GetMainInterfaceIpAddress() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AddBridgeInterface(t *testing.T) {
  type args struct {
    interfaceName string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AddBridgeInterface(tt.args.interfaceName)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AddBridgeInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AddBridgeInterface() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_ConfigureBridgeInterface(t *testing.T) {
  type args struct {
    interfaceName string
    ipAddr        string
    netmask       string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.ConfigureBridgeInterface(tt.args.interfaceName, tt.args.ipAddr, tt.args.netmask)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.ConfigureBridgeInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.ConfigureBridgeInterface() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_UpBridgeInterface(t *testing.T) {
  type args struct {
    interfaceName string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.UpBridgeInterface(tt.args.interfaceName)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.UpBridgeInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.UpBridgeInterface() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AddRouteToIp(t *testing.T) {
  type args struct {
    destIp    string
    gatewayIp string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AddRouteToIp(tt.args.destIp, tt.args.gatewayIp)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AddRouteToIp() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AddRouteToIp() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_DelRouteToIp(t *testing.T) {
  type args struct {
    destIp string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.DelRouteToIp(tt.args.destIp)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.DelRouteToIp() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.DelRouteToIp() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AddInterfaceToBridge(t *testing.T) {
  type args struct {
    bridgeInterfaceName string
    interfaceName       string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AddInterfaceToBridge(tt.args.bridgeInterfaceName, tt.args.interfaceName)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AddInterfaceToBridge() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AddInterfaceToBridge() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AddRerouteTrafficToGateway(t *testing.T) {
  type args struct {
    gatewayIp string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AddRerouteTrafficToGateway(tt.args.gatewayIp)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AddRerouteTrafficToGateway() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AddRerouteTrafficToGateway() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_RemoveRerouteTrafficToGateway(t *testing.T) {
  type args struct {
    gatewayIp string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.RemoveRerouteTrafficToGateway(tt.args.gatewayIp)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.RemoveRerouteTrafficToGateway() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.RemoveRerouteTrafficToGateway() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AllowIpForward(t *testing.T) {
  type args struct {
    on bool
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AllowIpForward(tt.args.on)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AllowIpForward() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AllowIpForward() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestCmdSuite_AllowTrafficForwardingOnInterface(t *testing.T) {
  type args struct {
    inputDevice  string
    outputDevice string
  }
  tests := []struct {
    name    string
    c       CmdSuite
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.c.AllowTrafficForwardingOnInterface(tt.args.inputDevice, tt.args.outputDevice)
      if (err != nil) != tt.wantErr {
        t.Errorf("CmdSuite.AllowTrafficForwardingOnInterface() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("CmdSuite.AllowTrafficForwardingOnInterface() = %v, want %v", got, tt.want)
      }
    })
  }
}
