package magent_base

import "testing"

func TestCreateConnectionToServicePeer(t *testing.T) {
  type args struct {
    peerPubKeyHash string
    keyFile        string
    isSecure       bool
    peerIp         string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := CreateConnectionToServicePeer(tt.args.peerPubKeyHash, tt.args.keyFile, tt.args.isSecure, tt.args.peerIp); (err != nil) != tt.wantErr {
        t.Errorf("CreateConnectionToServicePeer() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
