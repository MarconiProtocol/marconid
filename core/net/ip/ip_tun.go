package mnet_ip

import (
  "../../sys/cmd"
  "fmt"
  "time"
  "gitlab.neji.vm.tc/marconi/log"
)

func ConfigIpAddressForNewTunByCommand(
  taptunID string, ipAddr string, netmask string, peerIpAddr string, peerGatewayIpAddr string) (ret bool) {
  retCmdIpAddr := ConfigTunIpAddrByCommand(taptunID, ipAddr, netmask, peerIpAddr, peerGatewayIpAddr)
  fmt.Println("tun: ifconfig ", retCmdIpAddr, peerIpAddr)

  //TODO: call when peer is connected
  //NOTE: may not need to call when peer is connected/server side since this is tun which is for client or end point/grey
  if peerIpAddr != "" {
    gwIpAddr, _ := GetOwnGatewayIpAddress()
    //this needs to be called when client is connected or connection is made
    res, _ := ConfigRouteTargetIpAddr(peerIpAddr, gwIpAddr)

    //all traffic through Marconi
    //NOTE: during debug, disable this, sicne other end point will get all junk traffic from origin
    retRerouteAll, _ := RerouteAllTraffic(peerGatewayIpAddr)
    //retRerouteAll := "disabled for debugging"

    fmt.Println("route: brctl:white listing target ip address ",
      res, "gw ip:", gwIpAddr, "peer:", peerIpAddr, "reroute:", retRerouteAll)
  } else {
    fmt.Println("no peer ip address found or waiting mode - possible distribtion node mode/listening/waiting mode")

    //sysctl -w for system level forwarding enabling
    retTrafficForwardSys := allowTrafficForwardOnSystem()
    fmt.Println("retTrafficForwardSys: taptunID: ", retTrafficForwardSys)
    //TODO: add iptables/forarding on distribution node side
    //iptables -I FORWARD -i tun184 -o eth0 -j ACCEPT;
    //iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
    retTrafficForward := allowTrafficForwardOnTunInterface(taptunID)
    fmt.Println("retTrafficForward: taptunID: ", taptunID, "ret: ", retTrafficForward)
  }
  ret = true
  return
}

func ConfigTunIpAddrByCommand(
  taptunID string, ipAddr string, netmask string, peerIpAddr string, peerGatewayIpAddr string) (result map[int]string) {
  cmd, cmdArgs := GetCommandSetForAddIpAddressTunInterface(taptunID, ipAddr, netmask, peerIpAddr, peerGatewayIpAddr)
  fmt.Println("taptunID: ", taptunID, " ipAddr: ", ipAddr, " mask: ", netmask, "peer:", peerIpAddr)
  result = msys_cmd.ExecuteSequencialIdenticalCommand(cmd, cmdArgs)
  return
}

//L3 based = mconn = client mode
func ConfigMconnIpAddress(taptunNum string, ipAddr string, netmask string, gwIpAddr string, peerIpAddr string) {
  go func() {
    fmt.Println("ConfigBaseLayerIpAddress: sleep...: ", ipAddr, "netmask", netmask, "gw", gwIpAddr, "peer", peerIpAddr)
    time.Sleep(5 * time.Second)
    // NOTE: used gwIpAddr as peer gateway ip address for osx or client node's peer connection end
    retIpAddr := configClientLayerIpAddress(taptunNum, ipAddr, netmask, peerIpAddr, gwIpAddr)
    fmt.Println("ConfigBaseLayerIpAddress is called, ", retIpAddr, taptunNum, ipAddr, gwIpAddr)
    //TODO: mLocalIpGatewayIpAddr - route set
  }()
}

//3 based = client connenction
func configClientLayerIpAddress(taptunID string, ipAddr string, netmask string, peerIpAddr string, peerGatewayIpAddr string) bool {
  ret := ConfigIpAddressForNewTunByCommand(taptunID, ipAddr, netmask, peerIpAddr, peerGatewayIpAddr)
  fmt.Println("new new tap config with ip addr: ", ret)
  return ret
}

//gatewayIpAddr = internal marconi ip address of peer. this will be use to redirect traffic
func allowTrafficForwardOnTunInterface(taptunID string) (result map[int]string) {
  mlog.GetLogger().Debug("allowTrafficForwardOnTunInterface: taptunID: ", taptunID)
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    mlog.GetLogger().Fatal(err.Error())
  }
  // NOTE: replicating what was already written previously in ip_cmd,
  // which had these assumptions + hardcoding
  res, err := cmdSuite.AllowTrafficForwardingOnInterface("tun"+taptunID, "eth0")
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed to allow system ip forwarding using cmdSuite: %s", err.Error()))
  }
  return map[int]string{0: res}
}

func allowTrafficForwardOnSystem() (result map[int]string) {
  mlog.GetLogger().Debug("allowTrafficForwardOnSystem: called")
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    mlog.GetLogger().Fatal(err.Error())
  }
  res, err := cmdSuite.AllowIpForward(true)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed to allow system ip forwarding using cmdSuite: %s", err.Error()))
  }
  return map[int]string{0: res}
}
