package socks5

import (
  mnet_ip "../../../../core/net/ip"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "github.com/armon/go-socks5"
)

const SOCKS5_PORT string = "5656"

type SockServer struct {
  Server *socks5.Server
  ip     string
  port   string
}

/**
Serves as a temporary method by which our react native mobile app can connect to the server
Socks5Ip should be the mainInterface IP, and the port should be unique to the socks5 server
*/
func Initialize() (*SockServer, error) {
  sockServer := SockServer{}

  mainIpAddr, err := mnet_ip.GetMainInterfaceIpAddress()
  sockServer.ip = mainIpAddr

  sockServer.port = SOCKS5_PORT

  // Create a SOCKS5 server
  conf := &socks5.Config{}
  server, err := socks5.New(conf)
  if err != nil {
    return nil, err
  }
  sockServer.Server = server

  return &sockServer, nil
}

func (s *SockServer) Start() {
  mlog.GetLogger().Info(fmt.Sprintf("Starting socks5 server on %s:%s", s.ip, s.port))
  go func() {
    err := s.Server.ListenAndServe("tcp", s.ip+":"+s.port)
    if err != nil {
      mlog.GetLogger().Error(fmt.Sprintf("Socks Server failed! %s", err.Error()))
    }
  }()
}
