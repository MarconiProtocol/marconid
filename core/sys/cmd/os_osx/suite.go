package msys_cmd_osx

import (
  "../utils"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "os"
  "strings"
)

type CmdSuite string

func (c CmdSuite) GetOwnGatewayIpAddress() (string, error) {
  routeCmd, err := c.routeCmd(true)
  if err != nil {
    return "", err
  }
  grepCmd1, err := c.grepCmd("UG", false)
  if err != nil {
    return "", err
  }
  grepCmd2, err := c.grepCmd("UGH", true)
  if err != nil {
    return "", err
  }
  awkPrintCmd, err := c.awkPrintCmd(2)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunPipedCmds(routeCmd, grepCmd1, grepCmd2, awkPrintCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) GetMainInterfaceIpAddress() (string, error) {
  netstatCmd, err := c.netstatCmd(true, true)
  if err != nil {
    return "", err
  }
  grepCmd1, err := c.grepCmd("UG", false)
  if err != nil {
    return "", err
  }
  awkPrintCmd1, err := c.awkPrintCmd(6)
  if err != nil {
    return "", err
  }
  grepCmd2, err := c.grepCmd("inet", false)
  if err != nil {
    return "", err
  }
  awkPrintCmd2, err := c.awkPrintCmd(2)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunPipedCmds(netstatCmd, grepCmd1, awkPrintCmd1, grepCmd2, awkPrintCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddBridgeInterface(interfaceName string) (string, error) {
  return "", errors.New("AddBridgeInterface not supported on darwin")
}

func (c CmdSuite) ConfigureBridgeInterface(interfaceName string, ipAddr string, netmask string) (string, error) {
  return "", errors.New("ConfigureBridgeInterface not supported on darwin")
}

func (c CmdSuite) UpBridgeInterface(interfaceName string) (string, error) {
  return "", errors.New("UpBridgeInterface not supported on darwin")
}

func (c CmdSuite) AddRouteToIp(destIp string, gatewayIp string) (string, error) {
  routeCmd, err := c.routeAddRouteCmd(destIp, gatewayIp, 32)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(routeCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) DelRouteToIp(destIp string) (string, error) {
  return "", errors.New("stub")
}

func (c CmdSuite) AddInterfaceToBridge(bridgeInterfaceName string, interfaceName string) (string, error) {
  return "", errors.New("AddInterfaceToBridge not supported on darwin")
}

func (c CmdSuite) AddNetflowMonitorToBridge(bridgeInterfaceName string, collectorIp string, collectorPort string, netflowDirectory string) error {
  softflowdCmd, err := c.softflowdCmd(bridgeInterfaceName, collectorIp, collectorPort, 300)
  if err != nil {
    return errors.New(fmt.Sprintf("Command creation %s failed with error %s", softflowdCmd, err))
  }
  softflowdCmd.Stdout = nil
  softflowdCmd.Stderr = nil
  err = softflowdCmd.Start()
  if err != nil {
    return errors.New(fmt.Sprintf("Command execution %s failed with error %s", softflowdCmd, err))
  }

  if (collectorIp == "127.0.0.1" || collectorIp == "localhost") && netflowDirectory != "" {
    mlog.GetLogger().Info("Netflow collector on localhost, starting nfcapd.")
    _, err := os.Stat(netflowDirectory)
    if err != nil {
      return errors.New(fmt.Sprintf("Netflow logging directory does not exist %s", netflowDirectory))
    }
    nfcapdCmd, err := c.nfcapdCmd(collectorIp, collectorPort, netflowDirectory)
    if err != nil {
      return errors.New(fmt.Sprintf("Command creation %s failed with error %s", softflowdCmd, err))
    }
    nfcapdCmd.Stdout = nil
    nfcapdCmd.Stderr = nil
    err = nfcapdCmd.Start()
    if err != nil {
      return errors.New(fmt.Sprintf("Command execution %s failed with error %s", softflowdCmd, err))
    }
  }
  return nil
}

func (c CmdSuite) AddRerouteTrafficToGateway(gatewayIp string) (string, error) {
  routeCmd1, err := c.routeAddRouteCmd("0.0.0.0", gatewayIp, 1)
  if err != nil {
    return "", err
  }
  routeCmd2, err := c.routeAddRouteCmd("128.0.0.0", gatewayIp, 1)
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(routeCmd1)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(routeCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) RemoveRerouteTrafficToGateway(gatewayIp string) (string, error) {
  routeCmd1, err := c.routeDelRouteCmd("0.0.0.0", gatewayIp, 1)
  if err != nil {
    return "", err
  }
  routeCmd2, err := c.routeDelRouteCmd("128.0.0.0", gatewayIp, 1)
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(routeCmd1)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(routeCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) AllowIpForward(on bool) (string, error) {
  return "", errors.New("AllowIpForward not supported on darwin")
}

func (c CmdSuite) AllowTrafficForwardingOnInterface(inputDevice string, outputDevice string) (string, error) {
  return "", errors.New("AllowTrafficForwardingOnInterface not supported on darwin")
}
