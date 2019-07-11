package mnet_ip

import (
  "../vars"
  "../../sys/cmd"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
)

/*
  Assign an ip address to a bridge network interface
 */
func ConfigBridgeIpAddrByCommand(bridgeID string, ipAddr string, netmask string) (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  interfaceName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeID
  res, err := cmdSuite.ConfigureBridgeInterface(interfaceName, ipAddr, netmask)
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to configure bridge [%s] with addr %s and netmask %s using cmdSuite: %s", interfaceName, ipAddr, netmask, err.Error()))
  }
  return res, nil
}

/*
  "Up" the bridge interface
 */
func ConfigBridgeUpByCommand(bridgeID string) (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  // Bring the bridge interface up
  interfaceName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeID
  res, err := cmdSuite.UpBridgeInterface(interfaceName)
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to up the bridge interface [%s] using cmdSuite: %s", interfaceName, err.Error()))
  }
  return res, nil
}

/*
  Add a network interface to a bridge interface
 */
func ConfigTapToBridgeByCommand(bridgeID string, tapID string) (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }
  bridgeInterfaceName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeID
  interfaceName := mnet_vars.MPIPE_LINK_DEVICE_NAME_PREFIX + tapID
  res, err := cmdSuite.AddInterfaceToBridge(bridgeInterfaceName, interfaceName)
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to up add interface [%s] to bridge interface [%s] using cmdSuite: %s", bridgeInterfaceName, interfaceName, err.Error()))
  }
  return res, nil
}


/*
  Create and configure a bridge
 */
func ConfigBridgeByCommand(bridgeID string, ipAddr string, netmask string, resetBridge bool) (string, error) {
  cmdSuite, err := msys_cmd.GetSuite()
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to get cmd suite: %s", err))
  }

  // add bridge
  interfaceName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeID
  res, err := cmdSuite.AddBridgeInterface(interfaceName)
  if err != nil {
    return "", errors.New(fmt.Sprintf("Failed to create bridge [%s] using cmdSuite: %s", interfaceName, err.Error()))
  }
  mlog.GetLogger().Debug(fmt.Sprintf("Successfully created bridge [%s] using cmdSuite: %s", interfaceName, res))

  // configure bridge
  resultIpAddr, err := ConfigBridgeIpAddrByCommand(bridgeID, ipAddr, netmask)
  if err != nil {
    return "", nil
  }
  mlog.GetLogger().Debug(fmt.Sprintf("Successfully configured bridge [%s] with addr %s", interfaceName, ipAddr))

  return resultIpAddr, nil
}
