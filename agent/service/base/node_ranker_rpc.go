package magent_base

import (
  "../../../util"
  "net/http"
)

/*
  RPC to act as a health check for the node.
*/
func (agent *AgentClient) HandleRpcNodeFinderCheck() string {
  pkh, _ := mutil.GetInfohashByPubKey(&agent.edgeBeaconKey.PublicKey)
  // return our edgeBeaconKey, for confirmation purposes on other end
  return string([]byte(pkh))
}

type AgentClientService struct{}
type NodeFinderCheckArgs struct{}
type NodeFinderCheckReply struct {
  EdgeBeaconKey string
}

func (a *AgentClientService) NodeFinderCheckRPC(r *http.Request, args *NodeFinderCheckArgs, reply *NodeFinderCheckReply) error {
  reply.EdgeBeaconKey = Agent.HandleRpcNodeFinderCheck()
  return nil
}
