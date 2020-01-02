package rpc

import (
  mlog "github.com/MarconiProtocol/log"
  "github.com/gorilla/rpc/v2"
  "github.com/gorilla/rpc/v2/json2"
  "github.com/pkg/errors"
  "net/http"
)

const (
  RPC_PORT = "24802"
  RPC_PATH = "/rpc/m/request"
)

var server *rpc.Server

func init() {
  server = rpc.NewServer()
}

func StartRPCServer() {
  mlog.GetLogger().Infof("Starting JSON-RPC Server on port: %s", RPC_PORT)
  server.RegisterCodec(json2.NewCodec(), "application/json")
  http.Handle(RPC_PATH, server)
  mlog.GetLogger().Fatal(http.ListenAndServe(":"+RPC_PORT, nil))
}

func RegisterService(rcvr interface{}, name string) error {
  if server == nil {
    return errors.New("JSON-RPC Server has not been initialized")
  }
  return server.RegisterService(rcvr, name)
}
