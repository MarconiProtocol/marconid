package msys_cmd_ubuntu1804

import (
  "../../../../core"
  "../../../runtime"
  "../utils"
  "../vars"
  "bytes"
  "os/exec"
  "strings"
  "testing"
)

// todo: many operations in this test file need root privilege: sudo go test

var c CmdSuite = "cmdSuiteForTest"

func checkOSVersion(t *testing.T) {
  mruntime.GetMRuntime().SetRuntimeOS()
  if mruntime.GetMRuntime().GetRuntimeOS() != mcore.TYPE_OS_LINUX {
    t.Skip("Skip test")
  }
  distro, err := msys_cmd_utils.ParseLinuxVersion()
  if err != nil {
    t.Error("ParseLinuxVersion failed, err =", err)
  }
  if distro != msys_cmd_vars.UBUNTU1804 {
    t.Skip("Skip test")
  }
}

func TestCmdSuite_GetOwnGatewayIpAddress(t *testing.T) {
  checkOSVersion(t)
  // get own gateway ip address using bash script
  cmd := `route -n | grep UG | grep -v UGH | awk '{print $2}'`
  output, _ := exec.Command("sh", "-c", cmd).Output()
  ownGatewayIpAddress := strings.TrimSuffix(string(output), "\n")

  // check the result
  ip, err := c.GetOwnGatewayIpAddress()
  if err != nil || strings.Compare(ip, ownGatewayIpAddress) != 0 {
    t.Error("TestCmdSuite_GetOwnGatewayIpAddress failed")
  }
}

func TestCmdSuite_GetMainInterfaceIpAddress(t *testing.T) {
  checkOSVersion(t)
  // get main interface ip address using bash script
  cmd1 := `route -n | grep UG | grep -v UGH | awk '{print $8}'`
  output1, _ := exec.Command("sh", "-c", cmd1).Output()
  mainInterface := strings.TrimRight(string(output1), "\n")

  cmd2 := "ifconfig " + mainInterface + " | grep inet | cut -d ':' -f 2 | awk '{print $2}'"
  output2, _ := exec.Command("sh", "-c", cmd2).Output()
  // todo check the return value of GetMainInterfaceIpAddress, which has a redundant line feed
  mainInterfaceIpAddress := strings.Trim(string(output2), "\n")

  // check the result
  ip, err := c.GetMainInterfaceIpAddress()
  ip = strings.Trim(ip, "\n")
  if err != nil || strings.Compare(ip, mainInterfaceIpAddress) != 0 {
    t.Error("TestCmdSuite_GetMainInterfaceIpAddress failed")
  }
}

func TestCmdSuite_AddBridgeInterface(t *testing.T) {
  checkOSVersion(t)
  bridgeName := "testBridge"

  // add bridge
  _, err := c.AddBridgeInterface(bridgeName)
  if err != nil {
    t.Error("TestCmdSuite_AddBridgeInterface failed, err =", err)
  }

  // check if the bridge is added
  cmd1 := " brctl show | grep " + bridgeName + " | awk '{print $1}'"
  output1, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil || strings.Compare(bridgeName, strings.Trim(string(output1), "\n")) != 0 {
    t.Error("TestCmdSuite_AddBridgeInterface failed")
  }

  // clean up
  cmd2 := " brctl delbr " + bridgeName
  _, err = exec.Command("sh", "-c", cmd2).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddBridgeInterface failed, could not remove the added bridge interface, err =", err)
  }
}

func TestCmdSuite_ConfigureBridgeInterface(t *testing.T) {
  checkOSVersion(t)
  interfaceName := "testInterface"
  ipAddr := "10.27.16.12"
  netmask := "255.255.255.0"

  // create a dummy network interface and set it up
  cmd1 := "ip link add " + interfaceName + " type dummy && ip link set " + interfaceName + " up"
  _, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_ConfigureBridgeInterface failed, could not create the dummy network interface, err =", err)
  }

  // call ConfigureBridgeInterface
  _, err = c.ConfigureBridgeInterface(interfaceName, ipAddr, netmask)
  if err != nil {
    t.Error("TestCmdSuite_ConfigureBridgeInterface failed, ConfigureBridgeInterface returns err", err)
  }

  // get the ip address from ifconfig
  cmd2 := "ifconfig " + interfaceName + " | grep netmask | awk '{print $2}'"
  output2, err := exec.Command("sh", "-c", cmd2).Output()
  ipAddr2 := strings.Trim(string(output2), "\n")

  if err != nil {
    t.Error("TestCmdSuite_ConfigureBridgeInterface failed, could not get the IP address, err =", err)
  }

  // get the netmask from ifconfig
  cmd3 := "ifconfig " + interfaceName + " | grep netmask | awk '{print $4}'"
  output3, err := exec.Command("sh", "-c", cmd3).Output()
  netmask2 := strings.Trim(string(output3), "\n")
  if err != nil {
    t.Error("TestCmdSuite_ConfigureBridgeInterface failed, could not get the netmask, err =", err)
  }

  if ipAddr2 != ipAddr || netmask != netmask2 {
    t.Error("TestCmdSuite_ConfigureBridgeInterface failed, incorrect result, ip address =", ipAddr2, ", netmask =", netmask2)
  }

  // clean up
  cmd4 := "ip link delete " + interfaceName
  _, err = exec.Command("sh", "-c", cmd4).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not delete the dummy interface, err =", err)
  }
}

func TestCmdSuite_UpBridgeInterface(t *testing.T) {
  checkOSVersion(t)
  interfaceName := "testInterface"

  // create a dummy network interface
  cmd1 := "ip link add " + interfaceName + " type dummy"
  _, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_UpBridgeInterface failed, could not create the dummy network interface, err =", err)
  }

  // call UpBridgeInterface
  _, err = c.UpBridgeInterface(interfaceName)
  if err != nil {
    t.Error("TestCmdSuite_UpBridgeInterface failed, UpBridgeInterface returns error, err =", err)
  }

  // get the state of this interface
  //cmd2 := `ip link show | grep ` + interfaceName + ` | sed -n -e 's/^.*state //p' | awk '{print $1}'`
  cmd2 := `ifconfig -a | grep ` + interfaceName + ` | awk '{print $2}'`
  output2, err := exec.Command("sh", "-c", cmd2).Output()
  if err != nil {
    t.Error("TestCmdSuite_UpBridgeInterface failed, could not get the state of this interface, err =", err)
  }

  temp := strings.Trim(string(output2), "\n")
  index1 := strings.Index(temp, "<")
  index2 := strings.Index(temp, ",")
  state := temp[index1+1 : index2]

  if strings.Compare(state, "UP") != 0 {
    t.Error("TestCmdSuite_UpBridgeInterface failed, interface is not up, state =", state)
  }

  // clean up
  cmd3 := "ip link delete " + interfaceName
  _, err = exec.Command("sh", "-c", cmd3).Output()
  if err != nil {
    t.Error("TestCmdSuite_UpBridgeInterface failed, could not delete the dummy interface, err =", err)
  }
}

func TestCmdSuite_AddInterfaceToBridge(t *testing.T) {
  checkOSVersion(t)
  bridgeName := "testBridge"
  interfaceName := "testInterface"

  // create a dummy network interface and set it up
  cmd1 := "ip link add " + interfaceName + " type dummy && ip link set " + interfaceName + " up"
  _, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not create the dummy network interface, err =", err)
  }

  // add bridge
  _, err = c.AddBridgeInterface(bridgeName)
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not add the bridge, err =", err)
  }

  // check if the bridge is added
  cmd2 := "brctl show | grep " + bridgeName + " | awk '{print $1}'"
  output2, err := exec.Command("sh", "-c", cmd2).Output()
  if err != nil || strings.Compare(bridgeName, strings.Trim(string(output2), "\n")) != 0 {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, the bridge is not added")
  }

  // add interface to bridge
  _, err = c.AddInterfaceToBridge(bridgeName, interfaceName)
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not add the interface to the bridge, err =", err)
  }

  // check if the interface is added
  cmd3 := "brctl show " + bridgeName + " | grep " + interfaceName + " | awk '{print $4}'"
  output3, err := exec.Command("sh", "-c", cmd3).Output()
  if err != nil || strings.Compare(interfaceName, strings.Trim(string(output3), "\n")) != 0 {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, the interface is not added to the bridge")
  }

  // clean up
  cmd4 := "brctl delif " + bridgeName + " " + interfaceName
  _, err = exec.Command("sh", "-c", cmd4).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not delete the added interface, err =", err)
  }

  cmd5 := " brctl delbr " + bridgeName
  _, err = exec.Command("sh", "-c", cmd5).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not delete the added bridge, err =", err)
  }

  cmd6 := "ip link delete " + interfaceName
  _, err = exec.Command("sh", "-c", cmd6).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddInterfaceToBridge failed, could not delete the dummy interface, err =", err)
  }
}

func TestCmdSuite_AddRouteToIp(t *testing.T) {
  checkOSVersion(t)
  destIP := "172.217.0.46"

  // get the default gateway IP
  cmd1 := "ip route | grep default | awk '{print $3}'"
  output1, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddRouteToIp failed, could not get the default gateway, err =", err)
  }
  gatewayIP := strings.Trim(string(output1), "\n")

  // add route to IP
  _, err = c.AddRouteToIp(destIP, gatewayIP)
  if err != nil {
    t.Error("TestCmdSuite_AddRouteToIp failed, AddRouteToIp returns error:", err)
  }

  // check if new route is added
  cmd2 := "route -n | grep " + destIP
  output2, err := exec.Command("sh", "-c", cmd2).Output()
  splits := strings.Fields(string(output2))
  destIP2 := splits[0]
  gatewayIP2 := splits[1]
  netmask := splits[2]

  if destIP != destIP2 || gatewayIP != gatewayIP2 || netmask != "255.255.255.255" {
    t.Error("TestCmdSuite_AddRouteToIp failed, route is not added")
  }

  // clean up
  cmd3 := "route del " + destIP
  _, err = exec.Command("sh", "-c", cmd3).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddRouteToIp failed, could not delete the added route, err =", err)
  }
}

func TestCmdSuite_DelRouteToIp(t *testing.T) {
  checkOSVersion(t)
  destIP := "172.217.0.46"

  // get the default gateway IP
  cmd1 := "ip route | grep default | awk '{print $3}'"
  output1, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_DelRouteToIp failed, could not get the default gateway, err =", err)
  }
  gatewayIP := strings.Trim(string(output1), "\n")

  // add route to IP
  _, err = c.AddRouteToIp(destIP, gatewayIP)
  if err != nil {
    t.Error("TestCmdSuite_DelRouteToIp failed, AddRouteToIp returns error:", err)
  }

  _, err = c.DelRouteToIp(destIP)
  if err != nil {
    t.Error("TestCmdSuite_DelRouteToIp failed, DelRouteToIp returns error:", err)
  }
}

func TestCmdSuite_AddAndRemoveRerouteTrafficToGateway(t *testing.T) {
  checkOSVersion(t)
  gatewayIP := "127.0.0.1"

  // get the original kernel IP routing table
  cmd := "route -n"
  output1, err := exec.Command("sh", "-c", cmd).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, could not get the kernel IP routing table, err =", err)
  }

  // call AddRerouteTrafficToGateway
  _, err = c.AddRerouteTrafficToGateway(gatewayIP)
  if err != nil {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, AddRerouteTrafficToGateway returns error:", err)
  }

  // check if the kernel IP routing table is changed
  output2, err := exec.Command("sh", "-c", cmd).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, could not get the kernel IP routing table, err =", err)
  }

  if bytes.Compare(output1, output2) == 0 {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, did not add new item in kernel IP routing table")
  }

  // call RemoveRerouteTrafficToGateway
  _, err = c.RemoveRerouteTrafficToGateway(gatewayIP)
  if err != nil {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, AddRerouteTrafficToGateway returns error:", err)
  }

  // check kernel IP routing table again
  output3, err := exec.Command("sh", "-c", cmd).Output()
  if err != nil {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, could not get the kernel IP routing table, err =", err)
  }

  if bytes.Compare(output1, output3) != 0 {
    t.Error("TestCmdSuite_AddAndRemoveRerouteTrafficToGateway failed, could not recover the kernel IP routing table, err =", err)
  }
}

func TestCmdSuite_AllowIpForward(t *testing.T) {
  checkOSVersion(t)
  // get the original setting
  cmd1 := "sysctl net.ipv4.ip_forward | awk '{print $3}'"
  output1, err := exec.Command("sh", "-c", cmd1).Output()
  if err != nil {
    t.Error("TestCmdSuite_AllowIpForward failed, get the original setting failed, err =", err)
  }
  ipForward := strings.Trim(string(output1), "\n") == "1"

  // call AllowIpForward
  if ipForward {
    _, err = c.AllowIpForward(false)
  } else {
    _, err = c.AllowIpForward(true)
  }

  if err != nil {
    t.Error("TestCmdSuite_AllowIpForward failed, AllowIpForward returns error:", err)
  }

  // check if net.ipv4.ip_forward is set
  cmd2 := "sysctl net.ipv4.ip_forward | awk '{print $3}'"
  output2, err := exec.Command("sh", "-c", cmd2).Output()
  if err != nil {
    t.Error("TestCmdSuite_AllowIpForward failed, get the original setting failed, err =", err)
  }

  ipForward2 := strings.Trim(string(output2), "\n") == "1"

  if ipForward2 == ipForward {
    t.Error("TestCmdSuite_AllowIpForward failed, net.ipv4.ip_forward has not been changed")
  }

  // clean up
  if ipForward {
    _, err = c.AllowIpForward(true)
  } else {
    _, err = c.AllowIpForward(false)
  }
  if err != nil {
    t.Error("TestCmdSuite_AllowIpForward failed, cannot recover the original setting, err =", err)
  }
}

func TestCmdSuite_AllowTrafficForwardingOnInterface(t *testing.T) {
  // TODO test it later
}
