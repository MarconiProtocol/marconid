package mnet_ip

import (
  "../../../core"
  "../../runtime"
  "../../sys/cmd"
  "errors"
  "fmt"
  "strings"
)

/*
  Returns the main network interface ip address
*/
func GetMainInterfaceIpAddress() (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  mainInterfaceIp, err := cmdSuite.GetMainInterfaceIpAddress()
  if err != nil {
    return "", errors.New(fmt.Sprintf("cmdSuite.GetMainInterfaceIpAddress() failed with error: %s", err))
  }
  return strings.TrimSpace(mainInterfaceIp), nil
}

/*
  Returns the node's gateway ip address
*/
func GetOwnGatewayIpAddress() (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  gateway, err := cmdSuite.GetOwnGatewayIpAddress()
  if err != nil {
    return "", errors.New(fmt.Sprintf("cmdSuite.GetOwnGatewayIpAddress() failed with error: %s", err))
  }
  return gateway, nil
}

/////// Client Node L3 connection
// CMD changes based on platform.
// gatewayIpAddr is peer internal ip as gateway ip from client side connection
func GetCommandSetForAddIpAddressTunInterface(taptunID string, ipAddr string, netmask string, peerIpAddr string, gatewayIpAddr string) (string, map[int][]string) {
  var cmdArgs map[int][]string
  osType := mruntime.GetMRuntime().GetRuntimeOS()
  var cmd string
  if osType == mcore.TYPE_OS_LINUX {
    cmd = "ip"
    cmdArgs = map[int][]string{
      //0 : []string{"br" + bridgeID, "down"},
      0: []string{"tunnel", "add", "tun" + taptunID, "mode", "sit", "remote", ipAddr, "local", ipAddr},
      1: []string{"link", "set", "dev", "tun" + taptunID, "mtu", "1400", "netmask", netmask},
      2: []string{"link", "set", "dev", "tun" + taptunID, "up"},
    }
  } else if osType == mcore.TYPE_OS_DARWIN {
    cmd = "ifconfig"
    cmdArgs = map[int][]string{
      //0 : []string{"br" + bridgeID, "down"},
      //TODO: read below
      //utun0 for now or get it from return and set/override value. osx will not allow to pick or research more
      //https://developer.apple.com/legacy/library/documentation/Darwin/Reference/ManPages/man8/ifconfig.8.html
      // Above link no longer exists, closest replication https://www.freebsd.org/cgi/man.cgi?ifconfig(8)
      0: []string{"utun" + "0", "inet", ipAddr, gatewayIpAddr, "mtu", "1400"},
      1: []string{"utun" + "0", "up"},
    }
  } else if osType == mcore.TYPE_OS_WINDOWS {
    cmd = "netsh"
    cmdArgs = map[int][]string{
      0: []string{"interface" + "ip", "set", "address", ipAddr, gatewayIpAddr, "mtu", "1400"},
      1: []string{"utun" + "0", "up"},
    }
    //netsh interface ip set address "connection name" static 192.168.0.101 255.255.255.0 192.168.0.1
    //netsh interface ipv4 set subinterface “Local Area Connection” mtu=1458 store=persistent
  }
  return cmd, cmdArgs
}
