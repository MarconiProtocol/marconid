package msys_cmd_ubuntu1804

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
  ipRouteCmd, err := c.ipRouteCmd()
  if err != nil {
    return "", err
  }

  grepCmd1, err := c.grepCmd("default", false)
  if err != nil {
    return "", err
  }

  awkPrintCmd, err := c.awkPrintCmd(3)
  if err != nil {
    return "", err
  }

  res, err := msys_cmd_utils.RunPipedCmds(ipRouteCmd, grepCmd1, awkPrintCmd)
  if err != nil {
    return "", err
  }

  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) GetMainInterfaceIpAddress() (string, error) {
  ipRouteCmd, err := c.ipRouteCmd()
  if err != nil {
    return "", err
  }
  grepCmd1, err := c.grepCmd("default", false)
  if err != nil {
    return "", err
  }
  awkPrintCmd, err := c.awkPrintCmd(5)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunPipedCmds(ipRouteCmd, grepCmd1, awkPrintCmd)
  if err != nil {
    return "", err
  }

  // take the outStr from the previous commands and use as input for the next set
  // avoids having xargs
  ipCmd, err := c.ipCmd(strings.TrimSuffix(res, "\n"))
  if err != nil {
    return "", err
  }
  grepCmd3, err := c.grepCmd("inet", false)
  if err != nil {
    return "", err
  }
  grepCmd4, err := c.grepCmd("inet6", true)
  if err != nil {
    return "", err
  }
  cutCmd, err := c.cutCmd("/", "1")
  if err != nil {
    return "", err
  }
  awkPrintCmd2, err := c.awkPrintCmd(2)
  if err != nil {
    return "", err
  }
  res, err = msys_cmd_utils.RunPipedCmds(ipCmd, grepCmd3, grepCmd4, cutCmd, awkPrintCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddBridgeInterface(interfaceName string) (string, error) {
  brctlAddBrCmd, err := c.brctlAddBrCmd(interfaceName)
  if err != nil {
    return "", err
  }
  brctlStpCmd, err := c.brctlStpCmd(interfaceName, true)
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(brctlAddBrCmd)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(brctlStpCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) ConfigureBridgeInterface(interfaceName string, ipAddr string, netmask string) (string, error) {
  ipConfigureAddrCmd, err := c.ipConfigureAddrCmd(interfaceName, ipAddr, netmask)
  if err != nil {
    return "", err
  }
  ipUpCmd, err := c.ipUp(interfaceName, true)
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(ipConfigureAddrCmd)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(ipUpCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) UpBridgeInterface(interfaceName string) (string, error) {
  ipUpCmd, err := c.ipUp(interfaceName, true)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(ipUpCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddRouteToIp(destIp string, gatewayIp string) (string, error) {
  ipRouteCmd, err := c.ipAddRouteCmd(destIp, gatewayIp, 32, "")
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(ipRouteCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) DelRouteToIp(destIp string) (string, error) {
  ipCmd, err := c.ipDelRouteSimpleCmd(destIp)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(ipCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AddInterfaceToBridge(bridgeInterfaceName string, interfaceName string) (string, error) {
  brctlCmd, err := c.brctlAddIfCmd(bridgeInterfaceName, interfaceName)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(brctlCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
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
  ipCmd1, err := c.ipAddRouteCmd("0.0.0.0", gatewayIp, 1, "")
  if err != nil {
    return "", err
  }
  ipCmd2, err := c.ipAddRouteCmd("128.0.0.0", gatewayIp, 1, "")
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(ipCmd1)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(ipCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) RemoveRerouteTrafficToGateway(gatewayIp string) (string, error) {
  ipCmd1, err := c.ipDelRouteCmd("0.0.0.0", gatewayIp, 1, "")
  if err != nil {
    return "", err
  }
  ipCmd2, err := c.ipDelRouteCmd("128.0.0.0", gatewayIp, 1, "")
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(ipCmd1)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(ipCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}

func (c CmdSuite) AllowIpForward(on bool) (string, error) {
  sysctlCmd, err := c.sysctlSetIpForward(on)
  if err != nil {
    return "", err
  }
  res, err := msys_cmd_utils.RunCmd(sysctlCmd)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(res, "\n"), nil
}

func (c CmdSuite) AllowTrafficForwardingOnInterface(inputDevice string, outputDevice string) (string, error) {
  iptablesCmd1, err := c.iptablesInsert("", "FORWARD", inputDevice, outputDevice, "ACCEPT")
  if err != nil {
    return "", err
  }
  iptablesCmd2, err := c.iptablesInsert("nat", "POSTROUTING", "", outputDevice, "MASQUERADE")
  if err != nil {
    return "", err
  }
  res1, err := msys_cmd_utils.RunCmd(iptablesCmd1)
  if err != nil {
    return "", err
  }
  res2, err := msys_cmd_utils.RunCmd(iptablesCmd2)
  if err != nil {
    return "", err
  }
  return strings.TrimSuffix(fmt.Sprintf("%s\n%s", res1, res2), "\n"), nil
}
