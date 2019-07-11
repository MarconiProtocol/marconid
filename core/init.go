package mcore

import (
  "./config"
  "gitlab.neji.vm.tc/marconi/log"
)

func Initialize() {
  // initialize configuration
  mconfig.InitializeConfigs()
  // initialize logger
  mlog.Init(mconfig.GetAppConfig().Log.Dir, mconfig.GetAppConfig().Log.Level)
}
