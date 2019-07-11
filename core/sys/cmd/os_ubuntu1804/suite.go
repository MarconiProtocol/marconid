package msys_cmd_ubuntu1804

import (
	"fmt"
	"strings"
	"../utils"
)

type CmdSuite string

func (c CmdSuite) GetOwnGatewayIpAddress() (string, error) {
	routeCmd, err		:= c.routeCmd(true)
	grepCmd1, err 		:= c.grepCmd("UG", false)
	grepCmd2, err 		:= c.grepCmd("UGH", true)
	awkPrintCmd, err 	:= c.awkPrintCmd(2)
	outStr, err := msys_cmd_utils.RunPipedCmds(routeCmd, grepCmd1, grepCmd2, awkPrintCmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(outStr, "\n"), nil
}

func (c CmdSuite) GetMainInterfaceIpAddress() (string, error) {
	routeCmd, err 		:= c.routeCmd(true)
	if err != nil { return "", err }
	grepCmd1, err 		:= c.grepCmd("UG", false)
	if err != nil { return "", err }
	grepCmd2, err 		:= c.grepCmd("UGH", true)
	if err != nil { return "", err }
	awkPrintCmd, err 	:= c.awkPrintCmd(8)
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunPipedCmds(routeCmd, grepCmd1, grepCmd2, awkPrintCmd)
	if err != nil { return "", err }

	// take the outStr from the previous commands and use as input for the next set
	// avoids having xargs
	ifconfigCmd, err 	:= c.ifconfigCmd(strings.TrimSuffix(res, "\n"))
	if err != nil { return "", err }
	grepCmd3, err 		:= c.grepCmd("inet", false)
	if err != nil { return "", err }
	cutCmd, err 		:= c.cutCmd(":", "2")
	if err != nil { return "", err }
	awkPrintCmd2, err 	:= c.awkPrintCmd(2)
	if err != nil { return "", err }
	res, err = msys_cmd_utils.RunPipedCmds(ifconfigCmd, grepCmd3, cutCmd, awkPrintCmd2)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddBridgeInterface(interfaceName string) (string, error) {
	brctlAddBrCmd, err 	:= c.brctlAddBrCmd(interfaceName)
	if err != nil { return "", err }
	brctlStpCmd, err 	:= c.brctlStpCmd(interfaceName, true)
	if err != nil { return "", err }
	res1, err := msys_cmd_utils.RunCmd(brctlAddBrCmd)
	if err != nil { return "", err }
	res2, err := msys_cmd_utils.RunCmd(brctlStpCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) ConfigureBridgeInterface(interfaceName string, ipAddr string, netmask string) (string, error) {
	ifconfigAddrCmd, err := c.ifconfigIpAddrCmd(interfaceName, ipAddr, netmask)
	if err != nil { return "", err }
	ifconfigUpCmd, err := c.ifconfigUp(interfaceName, true)
	if err != nil { return "", err }
	res1, err := msys_cmd_utils.RunCmd(ifconfigAddrCmd)
	if err != nil { return "", err }
	res2, err := msys_cmd_utils.RunCmd(ifconfigUpCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) UpBridgeInterface (interfaceName string) (string, error) {
	ifconfigUpCmd, err := c.ifconfigUp(interfaceName, true)
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunCmd(ifconfigUpCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddRouteToIp(destIp string, gatewayIp string) (string, error) {
	routeCmd, err := c.routeAddRouteCmd(destIp, gatewayIp, "255.255.255.255", "")
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunCmd(routeCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) DelRouteToIp(destIp string) (string, error) {
	routeCmd, err := c.routeDelRouteSimpleCmd(destIp)
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunCmd(routeCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddInterfaceToBridge(bridgeInterfaceName string, interfaceName string) (string, error) {
	brctlCmd, err := c.brctlAddIfCmd(bridgeInterfaceName, interfaceName)
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunCmd(brctlCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddRerouteTrafficToGateway(gatewayIp string) (string, error) {
	routeCmd1, err := c.routeAddRouteCmd("0.0.0.0", gatewayIp, "128.0.0.0", "")
	if err != nil { return "", err }
	routeCmd2, err := c.routeAddRouteCmd("128.0.0.0", gatewayIp, "128.0.0.0", "")
	if err != nil { return "", err }
	res1, err := msys_cmd_utils.RunCmd(routeCmd1)
	if err != nil { return "", err }
	res2, err := msys_cmd_utils.RunCmd(routeCmd2)
	if err != nil { return "", err }
	return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) RemoveRerouteTrafficToGateway(gatewayIp string) (string, error) {
	routeCmd1, err := c.routeDelRouteCmd("0.0.0.0", gatewayIp, "128.0.0.0", "")
	if err != nil { return "", err }
	routeCmd2, err := c.routeDelRouteCmd("128.0.0.0", gatewayIp, "128.0.0.0", "")
	if err != nil { return "", err }
	res1, err := msys_cmd_utils.RunCmd(routeCmd1)
	if err != nil { return "", err }
	res2, err := msys_cmd_utils.RunCmd(routeCmd2)
	if err != nil { return "", err }
	return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) AllowIpForward(on bool) (string, error) {
	sysctlCmd, err := c.sysctlSetIpForward(on)
	if err != nil { return "", err }
	res, err := msys_cmd_utils.RunCmd(sysctlCmd)
	if err != nil { return "", err }
	return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AllowTrafficForwardingOnInterface(inputDevice string, outputDevice string) (string, error) {
	iptablesCmd1, err := c.iptablesInsert("", "FORWARD", inputDevice, outputDevice, "ACCEPT")
	if err != nil { return "", err }
	iptablesCmd2, err := c.iptablesInsert("nat", "POSTROUTING", "", outputDevice, "MASQUERADE")
	if err != nil { return "", err }
	res1, err := msys_cmd_utils.RunCmd(iptablesCmd1)
	if err != nil { return "", err }
	res2, err := msys_cmd_utils.RunCmd(iptablesCmd2)
	if err != nil { return "", err }
	return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}
