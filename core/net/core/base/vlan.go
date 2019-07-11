// +build linux

package mnet_core_base

import (
  "fmt"
  "github.com/vishvananda/netlink/nl"
  "gitlab.neji.vm.tc/marconi/log"
  "gitlab.neji.vm.tc/marconi/netlink"
  "io/ioutil"
)

// Enables vlan filtering on a network link
func EnableVlanFiltering(linkName string) {
  if !isVlanFilteringOn(linkName) {
    mlog.GetLogger().Info("enabling vlan filtering on link: ", linkName)
    setVlanFiltering(linkName, "1")
  } else {
    mlog.GetLogger().Info("vlan filtering already enabled for link: ", linkName)
  }
}

// Disables vlan filtering on a network link
func DisableVlanFilering(linkName string) {
  if isVlanFilteringOn(linkName) {
    mlog.GetLogger().Info("disabling vlan filtering on link: ", linkName)
    setVlanFiltering(linkName, "0")
  } else {
    mlog.GetLogger().Info("vlan filtering already disabled for link: ", linkName)
  }
}

// Checks if vlan filter is on
func isVlanFilteringOn(linkName string) bool {
  vfd := fmt.Sprintf("/sys/devices/virtual/net/%s/bridge/vlan_filtering", linkName)
  isEnable, err := ioutil.ReadFile(vfd)
  if err == nil {
    return string(isEnable) == "1" || string(isEnable) == "1\n"
  } else {
    mlog.GetLogger().Error("isVlanFilteringOn failed: ", err)
  }
  return false
}

// Sets the on/of state of vlan filtering on a network link
func setVlanFiltering(linkName string, state string) bool {
  vfd := fmt.Sprintf("/sys/devices/virtual/net/%s/bridge/vlan_filtering", linkName)
  if err := ioutil.WriteFile(vfd, []byte(state), 0644); err != nil {
    mlog.GetLogger().Info("setVlanFiltering: vlan filtering state: " + state + " - link: " + linkName)
    mlog.GetLogger().Info("setVlanFiltering: " + vfd)
    return false
  }
  return true
}

// Adds a new vlan filter entry
// equivalent to command: `bridge vlan add dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func AddVlanFilter(linkName string, vid uint16, pvid bool, untagged bool, self bool, master bool) {
  mlog.GetLogger().Info("AddVlanFilter ", linkName+", vid=", vid, ", pvid=", pvid, ", untagged=", untagged, ", self=", self, ", master=", master)
  link := GetLink(linkName)
  //if isVlanFilteringOn(bridgeName) {
  err := netlink.BridgeVlanAdd(link, vid, pvid, untagged, self, master)
  if err != nil {
    mlog.GetLogger().Error("AddVlanFilter failed: ", err)
  }
  //} else {
  //  mlog.GetLogger().Info("Cannot add vlan to ", bridgeName, ": vlan filter is not enabled on the bridge")
  //}
}

// Deletes a new vlan filter entry
// equivalent to command: `bridge vlan del dev DEV vid VID [ pvid ] [ untagged ] [ self ] [ master ]`
func DelVlanFilter(linkName string, vid uint16, pvid bool, untagged bool, self bool, master bool) {
  mlog.GetLogger().Info("DelVlanFilter ", linkName+", vid=", vid, ", pvid=", pvid, ", untagged=", untagged, ", self=", self, ", master=", master)
  link := GetLink(linkName)
  //if isVlanFilteringOn(bridgeName) {
  err := netlink.BridgeVlanDel(link, vid, pvid, untagged, self, master)
  if err != nil {
    mlog.GetLogger().Error("DelVlanFilter failed: ", err)
  }
  //} else {
  //  mlog.GetLogger().Info("Cannot remove vlan from ", bridgeName, ": vlan filter is not enabled on the bridge")
  //}
}

// Returns a map of vlan infos
func GetVlanMap() map[int32][]*nl.BridgeVlanInfo {
  vlanInfo, err := netlink.BridgeVlanList()
  if err != nil {
    mlog.GetLogger().Error("Failed to get vlan info")
  }
  return vlanInfo
}

func SetVfVlan(linkName string, vf int, vlan int) {
  link := GetLink(linkName)
  err := netlink.LinkSetVfVlan(link, vf, vlan)
  if err != nil {
    mlog.GetLogger().Error("LinkSetVfVlan failed: ", err)
  }
}
