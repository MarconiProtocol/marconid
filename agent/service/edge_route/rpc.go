package magent_edge_route_provider

import (
  "../../../core/crypto/key"
  "../../../core/peer"
  "errors"
  "fmt"
  "github.com/MarconiProtocol/log"
  "net/http"
  "strconv"
)

/*
	RPC handler to handle when an edge node requests for an edge connection
*/

type EdgeConnectionService struct{}

type EdgeConnectionArgs struct {
  PeerPubKeyHash string
  Identifier     string
}

type EdgeConnectionReply struct {
  Port              string
  AssignedIPNetwork string
  Pubkeyhash        string
}

func (e *EdgeConnectionService) RequestEdgeConnectionService(r *http.Request, args *EdgeConnectionArgs, reply *EdgeConnectionReply) error {
  mlog.GetLogger().Info("RECEIVED RPC REQUEST EDGE CONNECTION")
  peerPubKeyHash := args.PeerPubKeyHash

  // Add peer as a edge peer
  peerManager := mpeer.PeerManagerInstance()
  peerManager.AddPeer(mpeer.EDGE_PEER, peerPubKeyHash)

  // basic test here, we just want to start the pipe
  requesterIp := r.RemoteAddr
  edgeConnection, err := EdgeConnectionManagerInstance().RequestEdgeConnection(requesterIp)
  if err != nil {
    mlog.GetLogger().Errorf("Failed to allocate resources for edgeConnection for requester: %s", requesterIp)
    return errors.New(fmt.Sprintf("Failed to allocate resources for edgeConnection, err=%v", err))
  }
  port := strconv.Itoa(int(edgeConnection.Port))
  assignedIPNetwork := strconv.Itoa(int(edgeConnection.AssignedIPNetwork))
  mlog.GetLogger().Info(fmt.Sprintf("allocated port %s and network %s", port, assignedIPNetwork))

  // write response
  pubkeyhash := mcrypto_key.KeyManagerInstance().GetBasePublicKeyHash()

  reply.AssignedIPNetwork = assignedIPNetwork
  reply.Port = port
  reply.Pubkeyhash = pubkeyhash

  // Start edge route
  //go StartEdgeRoute(peerPubKeyHash, edgeConnection.Port, identifier)
  go StartEdgeRoute(peerPubKeyHash, edgeConnection.Port, peerPubKeyHash, assignedIPNetwork)
  return nil
}
