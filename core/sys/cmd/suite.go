package msys_cmd

import (
  "../../../core"
  "../../runtime"
  "./os_centos7"
  "./os_osx"
  "./os_ubuntu1604"
  "./os_ubuntu1804"
  "./utils"
  "./vars"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "github.com/pkg/errors"
)

type CmdSuite struct {
  CmdSuiteInt
}

/*
	CmdSuiteInt encapsulates actions that rely on os specific cli tools
*/
type CmdSuiteInt interface {
  GetOwnGatewayIpAddress() (string, error)
  GetMainInterfaceIpAddress() (string, error)
  AddBridgeInterface(interfaceName string) (string, error)
  ConfigureBridgeInterface(interfaceName string, ipAddr string, netmask string) (string, error)
  UpBridgeInterface(interfaceName string) (string, error)
  AddRouteToIp(destIp string, gatewayIp string) (string, error)
  DelRouteToIp(destIp string) (string, error)
  AddInterfaceToBridge(bridgeInterfaceName string, interfaceName string) (string, error)
  AddNetflowMonitorToBridge(bridgeInterfaceName string, collectorIp string, collectorPort string, netflowDirectory string) error
  AddRerouteTrafficToGateway(gatewayIp string) (string, error)
  RemoveRerouteTrafficToGateway(gatewayIp string) (string, error)
  AllowIpForward(on bool) (string, error)
  AllowTrafficForwardingOnInterface(inputDevice string, outputDevice string) (string, error)
}

/*
	Factory method to return a CmdSuite with the appropriate os cmd suite
*/
func GetSuite() (*CmdSuite, error) {
  var suite CmdSuite

  switch mruntime.GetMRuntime().GetRuntimeOS() {
  case mcore.TYPE_OS_LINUX:
    // need to find out the flavor/distro for linux
    distro, err := msys_cmd_utils.ParseLinuxVersion()
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Failed to parse linux version: %s", err.Error()))
    }
    switch distro {
    case msys_cmd_vars.UBUNTU1604:
      mlog.GetLogger().Debug("RETURNING THE U1604 SUITE")
      var osCmdSuite msys_cmd_ubuntu1604.CmdSuite
      suite = CmdSuite{osCmdSuite}
    case msys_cmd_vars.UBUNTU1804:
      mlog.GetLogger().Debug("RETURNING THE U1804 SUITE")
      var osCmdSuite msys_cmd_ubuntu1804.CmdSuite
      suite = CmdSuite{osCmdSuite}
    case msys_cmd_vars.CENTOS7:
      mlog.GetLogger().Debug("RETURNING THE C7 SUITE")
      var osCmdSuite msys_cmd_centos7.CmdSuite
      suite = CmdSuite{osCmdSuite}
    default:
      return nil, errors.New("Could not find appropriate cmd suite for your os")
    }
  case mcore.TYPE_OS_DARWIN:
    mlog.GetLogger().Debug("RETURNING THE DARWIN SUITE")
    var osCmdSuite msys_cmd_osx.CmdSuite
    suite = CmdSuite{osCmdSuite}
  default:
    return nil, errors.New("Could not find appropriate cmd suite for your os")
  }

  return &suite, nil
}
