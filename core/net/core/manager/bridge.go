package mnet_core_manager

import (
  "../../../config"
  "../../../../util"
  "../../ip"
  "../../vars"
  "../base"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
  "gitlab.neji.vm.tc/marconi/netlink"
  "strconv"
)

/*
	Get an already allocated bridge settings, or allocate a new one and return it
*/
func (nm *NetCoreManager) GetOrAllocateBridgeInfoForNetwork(netType NetworkType, netId string) (*BridgeInfo, bool, error) {
  nm.Lock()
  defer nm.Unlock()
  var newlyAllocated bool

  switch netType {
  case SERVICE_NET:
    if _, exists := nm.netIDtoBridgeInfoMap[netId]; !exists {
      bridgeIdInt, err := nm.allocateNextServiceBridgeId()
      if err != nil {
        return nil, false, errors.New(fmt.Sprintf("Could not create a bridge Id for network: %s of type SERVICE_NET", netId))
      }
      newlyAllocated = true
      nm.netIDtoBridgeInfoMap[netId] = &BridgeInfo{
        ID: strconv.Itoa(int(bridgeIdInt)),
        IP: "",
      }
    }
    return nm.netIDtoBridgeInfoMap[netId], newlyAllocated, nil
  default:
    return nil, false, errors.New(fmt.Sprintf("NetType with value [%d] does not exist", netType))
  }
}

/*
	Return the bridge info for a specific network type and id
*/
func (nm *NetCoreManager) GetBridgeInfoForNetwork(netType NetworkType, netId string) (*BridgeInfo, error) {
  nm.Lock()
  defer nm.Unlock()

  switch netType {
  case SERVICE_NET:
    if _, exists := nm.netIDtoBridgeInfoMap[netId]; !exists {
      return nil, errors.New(fmt.Sprintf("Bridge info does not exist for network of type: %d, and id: %s", netType, netId))
    }
    return nm.netIDtoBridgeInfoMap[netId], nil
  default:
    return nil, errors.New(fmt.Sprintf("NetType with value [%d] does not exist", netType))
  }
}

/*
	Create a bridge interface
*/
func (nm *NetCoreManager) CreateBridge(bridgeInfo *BridgeInfo, ipAddr string, netmask int, resetBridge bool) {
  // NOTE: Netlink API calls may fail in container or VPS environments where containers are sharing the same system resource
  // Therefore we must have a fallback option of using CMD for the system calls like CreateBinder and CreatePipe.
  // Also, we have found that for the CreateBinder and CreatePipe calls, sometimes the calls would return error but the Bridge/Pipe are created
  // successfully. So the true check is using GetLink call to see if the bridge or pipe exists, using on the error message returns from the
  // system call is not reliable.
  var bridge *netlink.Bridge
  bridgeName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeInfo.ID
  bridgeLink := mnet_core_base.GetBridge(bridgeName)
  if bridgeLink != nil {
    mlog.GetLogger().Debug(fmt.Sprintf("bridge already exists %s skipping creation", bridgeName))
  } else {
    mlog.GetLogger().Info(fmt.Sprintf("bridge doesn't exist, creating bridge: %s", bridgeName))
    bridge = mnet_core_base.CreateBinder(bridgeInfo.ID)
    if bridge != nil {
      mlog.GetLogger().Debug(fmt.Sprintf("bridge created successfully with NetLink: %s", bridge.Name))
    } else {
      mlog.GetLogger().Debug("bridge creation failed with NetLink, attempting again with CMD")
      retCmdBridge, err := mnet_ip.ConfigBridgeByCommand(bridgeInfo.ID, ipAddr, mutil.Get32BitMaskFromCIDR(netmask), resetBridge)
      if err != nil {
        mlog.GetLogger().Error(fmt.Sprintf("Failed to create bridge created with CMD: %s", retCmdBridge))
      } else {
        mlog.GetLogger().Info(fmt.Sprintf("bridge created with CMD: %s", retCmdBridge))
      }
    }

    // enable vlan filter
    if mconfig.GetAppConfig().Global.VlanFilterEnabled {
      mnet_core_base.EnableVlanFiltering(bridgeName)
    }
  }

  // Note: unfortunately OperState is not implemented in Netlink yet due to Linux restricting the write access
  // to network interface files related to the States, for now we have to use ifconfig up to bring up the
  // bridge after calling CreateBinder.
  res, err := mnet_ip.ConfigBridgeUpByCommand(bridgeInfo.ID)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("Failed bring up bridge with id: %s", bridgeInfo.ID))
  } else {
    mlog.GetLogger().Debug(fmt.Sprintf("Successfully brought up the bridge with id %s, res: %s", bridgeInfo.ID, res))
  }

}

/*
	Removes all bridge interfaces
*/
func (nm *NetCoreManager) RemoveAllBridges() {
  nm.Lock()
  defer nm.Unlock()

  for netId, bridgeInfo := range nm.netIDtoBridgeInfoMap {
    mlog.GetLogger().Debug(fmt.Sprintf("Deleting bridge: %s for service network: %s", netId, bridgeInfo.ID))
    mnet_core_base.DeleteBinder(bridgeInfo.ID)
  }
}

/*
	Bind a connection to a bridge interface
*/
func (nm *NetCoreManager) AddConnectionToBridge(bridgeInfo *BridgeInfo, taptunID string, peerIpAddr string) error {
  bridgeName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeInfo.ID
  bridge := mnet_core_base.GetBridge(bridgeName)

  pipeName, err := mnet_core_base.CreatePipe(taptunID, mnet_core_base.PIPE_MTU)
  if err != nil {
    return err
  } else {
    if pipeName != "" && bridge != nil {
      mlog.GetLogger().Info("Adding pipe ", pipeName, " to bridge ", bridge.Name)
      mnet_core_base.AddPipeIntoBinder(bridge, mnet_core_base.GetLink(pipeName))
    } else {
      mlog.GetLogger().Warn("CreatePipe", taptunID, "failed with NetLink, re-trying with CMD")
      _, err := mnet_ip.ConfigTapToBridgeByCommand(bridgeInfo.ID, taptunID)
      if err != nil {
        return err
      }
    }
  }

  // add vlan tag to the pipe
  vid := 2 // TODO get vlan ID from middleware
  if mconfig.GetAppConfig().Global.VlanFilterEnabled {
    mnet_core_base.AddVlanFilter(pipeName, uint16(vid), true, true, false, true)
  }
  if peerIpAddr != "" {
    gwIpAddr, err := mnet_ip.GetOwnGatewayIpAddress()
    if err != nil {
      return err
    }
    //this needs to be called when client is connected or connection is made
    _, err = mnet_ip.ConfigRouteTargetIpAddr(peerIpAddr, gwIpAddr)
    if err != nil {
      return err
    }
    mlog.GetLogger().Info(fmt.Sprintf("Added a route to: %s with gw: %s", peerIpAddr, gwIpAddr))

  } else {
    mlog.GetLogger().Info("no peer IP address found")
  }

  return nil
}

/*
	Assign an IP to a bridge interface
*/
func (nm *NetCoreManager) AssignIpAddrToBridge(bridgeInfo *BridgeInfo, ipAddr string, netmask int) {
  ifName := mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + bridgeInfo.ID

  nm.Lock()
  defer nm.Unlock()

  if mnet_core_base.GetLink(ifName) == nil || bridgeInfo == nil {
    mlog.GetLogger().Info("Bridge", ifName, "does not exist, skipping AssignIpAddrToBridge")
    return
  }

  // update with new IP
  bridgeInfo.IP = ipAddr

  // panic can happens when the environment doesn't have access to the kernel calls (e.g. docker)
  // regain control of execution when it happens and try configuring the IP via cmd
  defer func() {
    if r := recover(); r != nil {
      mlog.GetLogger().Warn("AssignNetlinkIpAddress panic'ed, failed to config bridge IP. Re-trying with CMD", r)
      mnet_ip.ConfigBridgeIpAddrByCommand(bridgeInfo.ID, ipAddr, mutil.Get32BitMaskFromCIDR(netmask))
    }
  }()

  // an interface could be associated with a list of IPs, so we first remove all other IPs from the list.
  for _, address := range mnet_ip.GetNetlinkAllIpAddress(ifName, mnet_ip.FAMILY_ALL) {
    mnet_ip.RemoveNetlinkIpAddress(ifName, address.IP.String(), netmask)
  }
  // assign the new IP to the interface
  mnet_ip.AssignNetlinkIpAddress(ifName, ipAddr, netmask)
}
