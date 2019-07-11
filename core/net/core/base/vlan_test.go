package mnet_core_base

import (
  "os/exec"
  "testing"
)

func createBridge(bridgeName string) error {
  cmd := "brctl addbr " + bridgeName
  _, err := exec.Command("sh", "-c", cmd).Output()
  return err
}

func deleteBridge(bridgeName string) error {
  cmd := "brctl delbr " + bridgeName
  _, err := exec.Command("sh", "-c", cmd).Output()
  return err
}

func TestEnableAndDisableVlanFiltering(t *testing.T) {
  bridgeName := "testBridge"

  // create bridge
  err := createBridge(bridgeName)
  if err != nil {
    t.Error("TestEnableAndDisableVlanFiltering failed, could not create bridge, err =", err)
  }

  DisableVlanFilering(bridgeName)
  if isVlanFilteringOn(bridgeName) {
    t.Error("TestEnableAndDisableVlanFiltering failed, isVlanFilteringOn should not return true")
  }

  EnableVlanFiltering(bridgeName)
  if !isVlanFilteringOn(bridgeName) {
    t.Error("TestEnableAndDisableVlanFiltering failed, isVlanFilteringOn should not return false")
  }

  DisableVlanFilering(bridgeName)
  if isVlanFilteringOn(bridgeName) {
    t.Error("TestEnableAndDisableVlanFiltering failed, isVlanFilteringOn should not return true")
  }

  // clean up
  err = deleteBridge(bridgeName)
  if err != nil {
    t.Error("TestEnableAndDisableVlanFiltering failed, could not delete bridge, err =", err)
  }
}

func TestGetVlanMap(t *testing.T) {
  // TODO: Add test cases.
}

func TestAddAndDeleteVlanFilter(t *testing.T) {
  // TODO: Add test cases.
}
