package mcore

import (
  "./config"
  mlog "github.com/MarconiProtocol/log"
)

func Initialize(baseDir string) {
  // initialize configuration
  mconfig.InitializeConfigs(baseDir)
  // initialize logger
  mlog.Init(mconfig.GetAppConfig().Log.Dir, mconfig.GetAppConfig().Log.Level)
}
