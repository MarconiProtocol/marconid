package magent_middleware_interface

import (
  "reflect"
  "testing"
)

func TestSubscribeParams_GetParamsJsonString(t *testing.T) {
  tests := []struct {
    name string
    r    SubscribeParams
    want string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.r.GetParamsJsonString(); got != tt.want {
        t.Errorf("SubscribeParams.GetParamsJsonString() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestUnsubscribeParams_GetParamsJsonString(t *testing.T) {
  tests := []struct {
    name string
    r    UnsubscribeParams
    want string
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.r.GetParamsJsonString(); got != tt.want {
        t.Errorf("UnsubscribeParams.GetParamsJsonString() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_createJsonPayload(t *testing.T) {
  type args struct {
    method string
    params string
  }
  tests := []struct {
    name string
    args args
    want []byte
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := createJsonPayload(tt.args.method, tt.args.params); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("createJsonPayload() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestRegisterForPeerUpdates(t *testing.T) {
  type args struct {
    pubKeyHash string
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
      if err := RegisterForPeerUpdates(tt.args.pubKeyHash); (err != nil) != tt.wantErr {
        t.Errorf("RegisterForPeerUpdates() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestUnregisterForPeerUpdates(t *testing.T) {
  tests := []struct {
    name    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := UnregisterForPeerUpdates(); (err != nil) != tt.wantErr {
        t.Errorf("UnregisterForPeerUpdates() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func Test_sendJsonRpcOverHttp(t *testing.T) {
  type args struct {
    host             string
    port             string
    path             string
    jsonPayloadBytes []byte
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
      if err := sendJsonRpcOverHttp(tt.args.host, tt.args.port, tt.args.path, tt.args.jsonPayloadBytes); (err != nil) != tt.wantErr {
        t.Errorf("sendJsonRpcOverHttp() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
