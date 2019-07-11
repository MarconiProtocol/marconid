package mnet_ip

import (
  "../../sys/cmd"
  "errors"
  "fmt"
)

/*
  Add a route (routing table entry) to the targetIpAddr through the gatewayIpAddr
 */
func ConfigRouteTargetIpAddr(targetIpAddr string, gatewayIpAddr string) (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  res, _ := cmdSuite.DelRouteToIp(targetIpAddr)
  res, err = cmdSuite.AddRouteToIp(targetIpAddr, gatewayIpAddr)
  if err != nil {
    return "", errors.New(fmt.Sprintf("cmdSuite.AddRouteToIp(%s, %s) failed with error: %s", targetIpAddr, gatewayIpAddr, err))
  }

  return res, nil
}

/*
  Direct all net traffic to through a gateway
 */
func RerouteAllTraffic(gatewayIpAddr string) (string, error){
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  res, err := cmdSuite.AddRerouteTrafficToGateway(gatewayIpAddr)
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to reroute traffic to gateway [%s] using cmdSuite: %s", gatewayIpAddr, err.Error()))
  }
  return res, nil
}
