package magent_util

import (
  mnet_ip "../../../core/net/ip"
  "errors"
  mlog "github.com/MarconiProtocol/log"
  "net/http"
)

type NetflowService struct{}
type StartNetflowArgs struct {
  IP         string
  Port       string
  BridgeId   string
  LoggingDir string
}
type StartNetflowReply struct{}

/*
  Handler for incoming start netflow monitoring requests.
  Starts monitoring of interface
*/
func (n *NetflowService) StartNetflowRPC(r *http.Request, args *StartNetflowArgs, reply *StartNetflowReply) error {
  mlog.GetLogger().Debug("StartNetflow RPC received args", args)
  if args.LoggingDir != "" && args.IP != "127.0.0.1" && args.IP != "localhost" { // Logging dir and remote host
    return errors.New("RPC StartNetflow received invalid payload, cannot have logging dir on remote host")
  }

  if err := mnet_ip.NetflowMonitorBridge(args.BridgeId, args.IP, args.Port, args.LoggingDir); err != nil {
    mlog.GetLogger().Errorf("RPC_START_NETFLOW:: error writing to response, err: %s", err)
    return err
  }
  return nil
}
