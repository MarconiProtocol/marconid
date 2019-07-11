package mconfig

import (
  "fmt"
  "sync"
)

var appConfig *AppConfiguration
var appConfigLock sync.RWMutex

var userConfig *UserConfiguration
var userConfigLock sync.RWMutex

func GetAppConfig() *AppConfiguration {
  appConfigLock.RLock()
  defer appConfigLock.RUnlock()
  return appConfig
}

func GetUserConfig() *UserConfiguration {
  userConfigLock.RLock()
  defer userConfigLock.RUnlock()
  return userConfig
}

func LoadAppConfig() {
  newAppConfiguration := new(AppConfiguration)
  if err := appConfigViperInst.Unmarshal(newAppConfiguration); err != nil {
    fmt.Printf("unable to decode application config file into struct, %v", err)
  }
  appConfigLock.Lock()
  appConfig = newAppConfiguration
  appConfigLock.Unlock()
}

func LoadUserConfig() {
  newUserConfiguration := new(UserConfiguration)
  if err := userConfigViperInst.Unmarshal(newUserConfiguration); err != nil {
    fmt.Printf("unable to decode appConfig file into struct, %v", err)
  }
  userConfigLock.Lock()
  userConfig = newUserConfiguration
  userConfigLock.Unlock()
}
