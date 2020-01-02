package mnet_ip

import "testing"

func TestConfigRouteTargetIpAddr(t *testing.T) {
  type args struct {
    targetIpAddr  string
    gatewayIpAddr string
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
      got, err := ConfigRouteTargetIpAddr(tt.args.targetIpAddr, tt.args.gatewayIpAddr)
      if (err != nil) != tt.wantErr {
        t.Errorf("ConfigRouteTargetIpAddr() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ConfigRouteTargetIpAddr() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestRerouteAllTraffic(t *testing.T) {
  type args struct {
    gatewayIpAddr string
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
      got, err := RerouteAllTraffic(tt.args.gatewayIpAddr)
      if (err != nil) != tt.wantErr {
        t.Errorf("RerouteAllTraffic() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("RerouteAllTraffic() = %v, want %v", got, tt.want)
      }
    })
  }
}
