package mnet_vars

import (
  "bytes"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
)

type ConnectionArgs struct {
  L2KeyFile      string              // base encryption and validation key/l2 data, timestamp
  EncPayload     bool                // "secure/nosecure"
  LocalPort      string              // local/remote binding port
  RemoteIpAddr   string              // remote ip if not listening mode
  RemotePort     string              // remote ip port
  DataKey        *bytes.Buffer       // DataKey bytes, produced through DH exchange
  DataKeySignal  *chan *bytes.Buffer // channel used to notify when key changes
  PeerPubKeyHash string              // peer's pubkey hash
}

func (c *ConnectionArgs) DebugPrint() {
  mlog.GetLogger().Info(fmt.Sprintf("ConnectionArgs - RemoteIP: %s, RemotePort: %s, PeerPubKeyHash: %s, L2KeyFile: %s", c.RemoteIpAddr, c.RemotePort, c.PeerPubKeyHash, c.L2KeyFile))
}
