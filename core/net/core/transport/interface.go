package mnet_core_transport

import (
  "../../if"
  "net"
)

type Transport interface {
  ListenAndTransmit(localIpAddr string, localPort string, targetIpAddress string, targetPort string,
    tapConn *mnet_if.Interface, key []byte, dataKey *[]byte, isSecure bool, isTun bool,
    listenSignalChannel chan string, transmitSignalChannel chan string) (net.Conn, error)
}
