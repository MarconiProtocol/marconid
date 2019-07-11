package msys_cmd_centos7

import (
	"errors"
	"fmt"
	"os/exec"
)

const (
	GREP_CMD 		= "/usr/bin/grep"
	AWK_CMD 		= "/usr/bin/awk"
	CUT_CMD 		= "/usr/bin/cut"
	IFCONFIG_CMD 	= "/usr/sbin/ifconfig"
	ROUTE_CMD		= "/usr/sbin/route"
	NETSTAT_CMD		= "/usr/bin/netstat"
	BRCTL_CMD		= "/usr/sbin/brctl"
	SYSCTL_CMD		= "/usr/sbin/sysctl"
	IPTABLES_CMD	= "/usr/sbin/iptables"
)

/*
	Grep for a pattern
		pattern - the pattern to match input on
		invert - whether the match should be inverted
*/
func (c CmdSuite) grepCmd(pattern string, invert bool) (*exec.Cmd, error) {
	grepCmd := exec.Command(GREP_CMD)
	if invert { grepCmd.Args = append(grepCmd.Args, "-v") }
	grepCmd.Args = append(grepCmd.Args, pattern)
	return grepCmd, nil
}

/*
	Awk command executing '{ print $X $Y $Z... }' where X Y Z are integers defined by indices
		indices - integer indices starting from 1 of the fields to print, where 1 is the first field
*/
func (c CmdSuite) awkPrintCmd(indices ...int) (*exec.Cmd, error) {
	awkPrintCmd := exec.Command(AWK_CMD)
	printIndices := ""
	for _, index := range indices {
		printIndices += fmt.Sprintf("$%d ", index)
	}
	printArg := fmt.Sprintf("{print %s}", printIndices)
	awkPrintCmd.Args = append(awkPrintCmd.Args, printArg)
	return awkPrintCmd, nil
}

/*
	Cut command helps to remove sections from lines of text/files
		delim - the delimiting character that separates fields, if left blank, the -d flag is not set and TAB will be used
		fields - only keeps the fields listed
 */
func (c CmdSuite) cutCmd(delim string, fields string)	(*exec.Cmd, error) {
	cutCmd := exec.Command(CUT_CMD)
	if delim != "" { cutCmd.Args = append(cutCmd.Args, "-d" + delim) }
	if fields != "" { cutCmd.Args = append(cutCmd.Args, "-f" + fields)}
	return cutCmd, nil
}

/*
	ifconfig command to display all interfaces or the specified interface
		iface - the name of the interface to be displayed
 */
func (c CmdSuite) ifconfigCmd(iface string) (*exec.Cmd, error) {
	ifconfigCmd := exec.Command(IFCONFIG_CMD)
	if iface != "" { ifconfigCmd.Args = append(ifconfigCmd.Args, iface) }
	return ifconfigCmd, nil
}

/*
	ifconfig command to configure a specified interfaces ip address and netmask
		iface - the name of the interface to be modified
		ipaddr - the ip address to be set to the interface
		netmask - the netmask to be set to the interface
 */
func (c CmdSuite) ifconfigIpAddrCmd(iface string, ipAddr string, netmask string) (*exec.Cmd, error) {
	ifconfigCmd := exec.Command(IFCONFIG_CMD)
	if iface != "" { ifconfigCmd.Args = append(ifconfigCmd.Args, iface, ipAddr, "netmask", netmask) }
	return ifconfigCmd, nil
}

/*
	ifconfig command to bring up/down a specified interface
		iface - the name of the interface to be brought up/down
		up - true => up / false => down
 */
func (c CmdSuite) ifconfigUp(iface string, up bool) (*exec.Cmd, error) {
	ifconfigCmd := exec.Command(IFCONFIG_CMD)
	status := ""
	if up { status = "up" } else { status = "down"}
	if iface != "" { ifconfigCmd.Args = append(ifconfigCmd.Args, iface, status)}
	return ifconfigCmd, nil
}

/*
	route command with the optional use of n flag
		nflag - prints the ip address instead of the hostnames
 */
func (c CmdSuite) routeCmd(nflag bool) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if nflag { routeCmd.Args = append(routeCmd.Args, "-n") }
	return routeCmd, nil
}

/*
	route command to add new route to routing table
		destIp - the destination IP
		gatewayIp - the gateway machine that packets destined for destIp will be forwarded
		netmask - the netmask associated with provided destIp
		device - the network device/interface that will be used to flush the traffic
 */
func (c CmdSuite) routeAddRouteCmd(destIp string, gatewayIp string, netmask string, device string) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if destIp == "" || gatewayIp == "" || netmask == "" {
		return nil, errors.New("destIp, gatewayIp, netmask must be provided to routeAddRouteCmd")
	}
	routeCmd.Args = append(routeCmd.Args, "add", "-net", destIp, "netmask", netmask, "gw", gatewayIp)
	if device != "" { routeCmd.Args = append(routeCmd.Args, "dev", device) }
	return routeCmd, nil
}

/*
	route command to delete existing route from routing table
		destIp - the destination IP
 */
func (c CmdSuite) routeDelRouteSimpleCmd(destIp string) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if destIp == "" {
		return nil, errors.New("destIp must be provided to routeDelRouteCmd")
	}
	routeCmd.Args = append(routeCmd.Args, "del", destIp)
	return routeCmd, nil
}

/*
	route command to delete existing route from routing table
		destIp - the destination IP
		gatewayIp - the gateway machine that packets destined for destIp will be forwarded
		netmask - the netmask associated with provided destIp
		device - the network device/interface that will be used to flush the traffic
 */
func (c CmdSuite) routeDelRouteCmd(destIp string, gatewayIp string, netmask string, device string) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if destIp == "" || gatewayIp == "" || netmask == "" {
		return nil, errors.New("destIp, gatewayIp, netmask must be provided to routeDelRouteCmd")
	}
	routeCmd.Args = append(routeCmd.Args, "del", "-net", destIp, "netmask", netmask, "gw", gatewayIp)
	if device != "" { routeCmd.Args = append(routeCmd.Args, "dev", device) }
	return routeCmd, nil
}

/*
	netstat command with optional use of r and n flags
		rflag - prints the kernel routing tables
		nflag - prints the ip address instead of the hostnames
 */
func (c CmdSuite) netstatCmd(rflag bool, nflag bool) (*exec.Cmd, error) {
	netstatCmd := exec.Command(NETSTAT_CMD)
	if rflag { netstatCmd.Args = append(netstatCmd.Args, "-r") }
	if nflag { netstatCmd.Args = append(netstatCmd.Args, "-n") }
	return netstatCmd, nil
}

/*
	brctl addbr command creates a new bridge interface
		bridgeIntName - the name of the bridge interface that will be created
 */
func (c CmdSuite) brctlAddBrCmd(bridgeIntName string) (*exec.Cmd, error) {
	brctlCmd := exec.Command(BRCTL_CMD)
	if bridgeIntName != "" { brctlCmd.Args = append(brctlCmd.Args, "addbr", bridgeIntName)}
	return brctlCmd, nil
}

/*
	brctl stp command updates an existing bridge interface to turn on/off stp
		bridgeIntName - the name of the bridge interface that will be modified
		on - true => spanning tree protocol will be turned on
 */
func (c CmdSuite) brctlStpCmd(bridgeIntName string, on bool) (*exec.Cmd, error) {
	brctlCmd := exec.Command(BRCTL_CMD)
	if bridgeIntName != "" { brctlCmd.Args = append(brctlCmd.Args, "stp", bridgeIntName, "on")}
	return brctlCmd, nil
}

/*
	brctl addif command adds an interface to a specified bridge interface
		bridgeIntName - the name of the bridge interface that will be modified
		intName - the name of the interface that will be added to the bridge
 */
func (c CmdSuite) brctlAddIfCmd(bridgeIntName string, intName string) (*exec.Cmd, error) {
	brctlCmd := exec.Command(BRCTL_CMD)
	if bridgeIntName == "" || intName == "" {
		return nil, errors.New("brctlAddIfCmd: bridgeIntName and intName must be defined")
	}
	brctlCmd.Args = append(brctlCmd.Args, "addif", bridgeIntName, intName)
	return brctlCmd, nil
}

/*
	use sysctl to modify the kernel parameter 'net.ipv4.ip_forward'
		on - true => set the parameter net.ipv4.ip_forward to 1
 */
func (c CmdSuite) sysctlSetIpForward(on bool) (*exec.Cmd, error) {
	sysctlCmd := exec.Command(SYSCTL_CMD)
	arg := ""
	if on {
		arg = "net.ipv4.ip_forward=1"
	} else {
		arg = "net.ipv4.ip_forward=0"
	}
	sysctlCmd.Args = append(sysctlCmd.Args, "-w", arg)
	return sysctlCmd, nil
}

/*
	use iptables to insert a rule into a specific chain for a given table
		tableName - the table name, if left blank, defaults to filter table (by not using -t flag)
		chain - the target chain to insert the rule to
		input - the network interface that is the input
		output - the network interface that is the output
		jump - target for rule
 */
func (c CmdSuite) iptablesInsert(tableName string, chain string, input string, output string, jump string) (*exec.Cmd, error) {
	iptablesCmd := exec.Command(IPTABLES_CMD)
	if tableName != "" {
		iptablesCmd.Args = append (iptablesCmd.Args, "-t", tableName)
	}
	iptablesCmd.Args = append (iptablesCmd.Args, "-I", chain, "-j", jump)
	if input != "" {
		iptablesCmd.Args = append (iptablesCmd.Args, "-i", input)
	}
	if output != "" {
		iptablesCmd.Args = append (iptablesCmd.Args, "-o", output)
	}
	return iptablesCmd, nil
}
