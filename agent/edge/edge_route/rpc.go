package magent_edge_route_request

import (
  "../../../core/crypto/key"
  "../../../core/rpc"
  "../../service/edge_route"
  "errors"
  "fmt"
)

type edgeRouteRequestResponseDetails struct {
  host               string
  port               string
  edgeNodePubKeyHash string
}

// TODO: we need to also send some additional identifier based on
// end-user -> FeO -> BeO relationship since there can be multiple edge routes to one BeO based on end-user
func sendRequestEdgeRouteRPC(target string, pkhIdentifier string) (*edgeRouteRequestResponseDetails, error) {
  args := magent_edge_route_provider.EdgeConnectionArgs{
    PeerPubKeyHash: mcrypto_key.KeyManagerInstance().GetBasePublicKeyHash(),
    Identifier:     pkhIdentifier,
  }
  reply := magent_edge_route_provider.EdgeConnectionReply{}

  client := rpc.NewRPCClient(target)
  if err := client.Call("EdgeConnectionService.RequestEdgeConnectionService", &args, &reply); err != nil {
    return nil, errors.New(fmt.Sprintf("sendRequestEdgeRouteRpc - Nil/error response from: %s; response: %v", target, err))
  }

  edgeRouteDetails := edgeRouteRequestResponseDetails{
    host:               magent_edge_route_provider.EDGE_NETWORK_FIRST_OCTET + reply.AssignedIPNetwork + magent_edge_route_provider.EDGE_CLIENT_NETWORK_HOST,
    port:               reply.Port,
    edgeNodePubKeyHash: reply.Pubkeyhash,
  }
  return &edgeRouteDetails, nil
}
