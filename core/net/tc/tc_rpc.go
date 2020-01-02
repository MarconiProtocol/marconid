package tc

import "net/http"

type TrafficControlService struct{}

func (t *TrafficControlService) SetNetemRPC(r *http.Request, args *NetemArgs, reply *NetemReply) error {
  return SetNetem(args.InterfaceName, args.Delay, args.Loss, args.Duplicate, args.ReorderProb, args.CorruptProb)
}

func (t *TrafficControlService) SetTbfRPC(r *http.Request, args *TbfArgs, reply *TbfReply) error {
  return SetTbf(args.InterfaceName, args.Bandwidth, args.LatencyInMillis)
}

func (t *TrafficControlService) ResetRPC(r *http.Request, args *ResetArgs, reply *ResetReply) error {
  return DeleteAllQdisc(args.InterfaceName)
}
