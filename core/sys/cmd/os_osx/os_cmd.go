package msys_cmd_osx

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

const (
	GREP_CMD 		= "/usr/bin/grep"
	AWK_CMD 		= "/usr/bin/awk"
	CUT_CMD 		= "/usr/bin/cut"
	IFCONFIG_CMD 	= "/sbin/ifconfig"
	ROUTE_CMD		= "/sbin/route"
	NETSTAT_CMD		= "/usr/sbin/netstat"
)

// NOTE: OSX is not ready, we should revisit this when we are working on support for OSX
// TRANSLATION: MAY NOT WORK / UNTESTED

/*
	Grep for a pattern
		pattern - the pattern to match input on
		invert - whether the match should be inverted
*/
func (c CmdSuite)  grepCmd(pattern string, invert bool) (*exec.Cmd, error) {
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
func (c CmdSuite) routeAddRouteCmd(destIp string, gatewayIp string, netmaskLength int) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if destIp == "" || gatewayIp == "" || !(netmaskLength > 0 && netmaskLength <= 32) {
		return nil, errors.New("destIp, gatewayIp, netmaskLength must be provided to routeAddRouteCmd")
	}
	destIp = destIp + "/" + strconv.Itoa(netmaskLength)
	routeCmd.Args = append(routeCmd.Args, "-n", "add", "-net", destIp, gatewayIp)
	return routeCmd, nil
}

/*
	route command to delete existing route from routing table
		destIp - the destination IP
		gatewayIp - the gateway machine that packets destined for destIp will be forwarded
		netmask - the netmask associated with provided destIp
		device - the network device/interface that will be used to flush the traffic
 */
func (c CmdSuite) routeDelRouteCmd(destIp string, gatewayIp string, netmaskLength int) (*exec.Cmd, error) {
	routeCmd := exec.Command(ROUTE_CMD)
	if destIp == "" || gatewayIp == "" || !(netmaskLength > 0 && netmaskLength <= 32) {
		return nil, errors.New("destIp, gatewayIp, netmaskLength must be provided to routeAddRouteCmd")
	}
	destIp = destIp + "/" + strconv.Itoa(netmaskLength)
	routeCmd.Args = append(routeCmd.Args, "-n", "delete", "-net", destIp, gatewayIp)
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
