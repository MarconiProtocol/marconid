// +build linux

package mnet_core_base

import (
  "errors"
  "fmt"
  "io/ioutil"
  "sync"

  "../../../../core"
  "../../vars"

  "../../../runtime"
  "gitlab.neji.vm.tc/marconi/log"
  "gitlab.neji.vm.tc/marconi/netlink"
)

//binder could exist in two mode bridge or bond
//bond is used to amplify

var mutexBinder sync.Mutex

func CreateBinder(binderNum string) (bridge *netlink.Bridge) {
  if mruntime.GetMRuntime().GetRuntimeOS() == mcore.TYPE_OS_LINUX {
    link := netlink.NewLinkAttrs()
    link.Name = mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + binderNum
    // TODO: OperState is not respected by Netlink currently, need a workaround to support it.
    //link.OperState = netlink.OperUp
    attr := &netlink.Bridge{LinkAttrs: link}
    err := netlink.LinkAdd(attr)
    if err != nil {
      mlog.GetLogger().Errorf("CreateBinder could not add %s: %v", link.Name, err)
    } else {
      mutexBinder.Lock()
      defer mutexBinder.Unlock()

      EnableLinkSTP(link.Name)
      bridge = attr
    }
  } else {
    mlog.GetLogger().Warn("CreateBinder OS is not supported:", mruntime.GetMRuntime().GetRuntimeOS())
  }
  return bridge
}

/*
** Remove the binder specified by binderNum, used during graceful shutdown of Marconid
 */
func DeleteBinder(binderNum string) {
  mutexBinder.Lock()
  defer mutexBinder.Unlock()

  bridgeLink := GetLink(mnet_vars.MBINDER_LINK_DEVICE_NAME_PREFIX + binderNum)
  if bridgeLink != nil {
    mlog.GetLogger().Info("DeleteBinders removing", bridgeLink.Attrs().Name)
    netlink.LinkDel(bridgeLink)
  }
}

func AddPipeIntoBinder(binder *netlink.Bridge, pipe netlink.Link) bool {
  err := netlink.LinkSetMaster(pipe, binder)
  if err != nil {
    mlog.GetLogger().Errorf("could not add mpipe: %v", err)
    return false
  }
  return true
}

func RemovePipeFromBinder(pipe netlink.Link) bool {
  err := netlink.LinkSetNoMaster(pipe)
  if err == nil {
    mlog.GetLogger().Errorf("could not remove mpipe: %v", err)
    return false
  }
  return true
}

func EnableLinkSTP(linkName string) {
  if !isEnableLinkSTP(linkName) {
    mlog.GetLogger().Info("enabling stp on link:", linkName)
    setLinkSTP(linkName, "1")
  } else {
    mlog.GetLogger().Info("already stp enabled for link:", linkName)
  }
}

func DisableLinkSTP(linkName string) {
  if isEnableLinkSTP(linkName) {
    mlog.GetLogger().Info("disabling stp on link:", linkName)
    setLinkSTP(linkName, "0")
  } else {
    mlog.GetLogger().Info("already stp disabled for link", linkName)
  }
}

func setLinkSTP(linkName string, state string) bool {
  ///sys/devices/virtual/net/mb200/bridge/stp_state
  ///sys/class/net/mb200/bridge/stp_state
  vfd := fmt.Sprintf("/sys/devices/virtual/net/%s/bridge/stp_state", linkName)
  if err := ioutil.WriteFile(vfd, []byte(state), 0644); err != nil {
    mlog.GetLogger().Info("setLinkSTP: stp state: " + state + " - link: " + linkName)
    mlog.GetLogger().Info("setnLinkSTP: " + vfd)
    return false
  }
  return true
}

func isEnableLinkSTP(linkName string) bool {
  vfd := fmt.Sprintf("/sys/devices/virtual/net/%s/bridge/stp_state", linkName)
  if isEnable, err := ioutil.ReadFile(vfd); err != nil && string(isEnable) == "1" {
    return true
  }
  return false
}

//regular l2 link
func CreatePipe(pipeNum string, mtu int) (string, error) {
  if mruntime.GetMRuntime().GetRuntimeOS() == mcore.TYPE_OS_LINUX {
    linkAttrs := netlink.NewLinkAttrs()
    linkAttrs.Name = mnet_vars.MPIPE_LINK_DEVICE_NAME_PREFIX + pipeNum
    linkAttrs.MTU = mtu
    attr := &netlink.Tuntap{
      LinkAttrs: linkAttrs,
      Mode:      netlink.TUNTAP_MODE_TAP,
    }

    err := netlink.LinkAdd(attr)
    if err != nil {
      if GetLink(linkAttrs.Name) != nil {
        mlog.GetLogger().Info("CreatePipe succeeded in creating ", linkAttrs.Name)
        return linkAttrs.Name, nil
      }
      return "", errors.New(fmt.Sprintf("CreatePipe could not add %s, err: %s", linkAttrs.Name, err))
    }
    return linkAttrs.Name, nil
  }

  return "", errors.New(fmt.Sprintf("CreatePipe OS is not supported %d", mruntime.GetMRuntime().GetRuntimeOS()))
}

func GetLink(ifName string) netlink.Link {
  link, err := netlink.LinkByName(ifName)
  if err != nil {
    return nil
  }
  return link
}

// return a handle to the bridge if it exists
func GetBridge(bridgeName string) *netlink.Bridge {
  link := GetLink(bridgeName)
  if link != nil {
    return link.(*netlink.Bridge)
  }
  return nil
}
