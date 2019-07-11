package mconfig

import (
  "../net/vars"
  "fmt"
  "github.com/fsnotify/fsnotify"
  "github.com/spf13/viper"
  "gitlab.neji.vm.tc/marconi/log"
  "os"
  "path/filepath"
)

const CONFIG_PATH = "/opt/marconi/etc/marconid"
const CONFIG_EXT = "yml"

var appConfigViperInst  *viper.Viper
var userConfigViperInst *viper.Viper

func InitializeConfigs() {
  InitializeAppConfig()
  InitializeUserConfig()
}

// initialize the appConfig object by parsing the appConfig file each time it's updated
// config_manager.go is kept outside of package "appConfig" on purpose so that we can avoid cyclic reference
// for packages that need to access the appConfig object, yet also expect to be notified when appConfig file changes
func InitializeAppConfig() {
  appConfigViperInst = viper.New()

  configName := "config"
  configFilename := filepath.Join(CONFIG_PATH, "config_dev."+CONFIG_EXT)
  if _, err := os.Stat(configFilename); !os.IsNotExist(err) {
    configName = "config_dev"
  }
  appConfigViperInst.SetConfigName(configName)
  appConfigViperInst.SetConfigType(CONFIG_EXT)
  appConfigViperInst.AddConfigPath(CONFIG_PATH)

  // set default value for required configs
  setAppConfigDefaults()

  // read in the appConfig file
  readAndLoadAppConfig()

  // triggered on initial read and every time appConfig is modified
  // make sure any functions called within OnConfigChange are thread-safe
  appConfigViperInst.OnConfigChange(func(event fsnotify.Event) {
    // read in appConfig file
    readAndLoadAppConfig()
    // update the log level
    mlog.SetOutputLevel(GetAppConfig().Log.Level)
    fmt.Println("Config file changed:", event.Name)
  })
  // watch the appConfig for file change
  appConfigViperInst.WatchConfig()
}

// tell viper to read in appConfig file and parse it
func readAndLoadAppConfig() {
  if err := appConfigViperInst.ReadInConfig(); err != nil {
    switch err.(type) {
    case viper.ConfigFileNotFoundError:
      // pass if appConfig file does not exist
    default:
      mlog.GetLogger().Error(fmt.Sprintf("Error reading app config file, %s\n", err.Error()))
      os.Exit(1)
    }
  }
  // parse appConfig file into their corresponding struct
  LoadAppConfig()
}


func InitializeUserConfig() {
  userConfigViperInst = viper.New()

  configName := "user_config"
  userConfigViperInst.SetConfigName(configName)
  userConfigViperInst.SetConfigType(CONFIG_EXT)
  userConfigViperInst.AddConfigPath(CONFIG_PATH)

  // set default value for required configs
  setUserConfigDefaults()

  // read in the appConfig file
  readAndLoadUserConfig()

  // triggered on initial read and every time appConfig is modified
  // make sure any functions called within OnConfigChange are thread-safe
  userConfigViperInst.OnConfigChange(func(event fsnotify.Event) {
    // read in appConfig file
    readAndLoadUserConfig()
    fmt.Println("Config file changed:", event.Name)
  })
  // watch the appConfig for file change
  userConfigViperInst.WatchConfig()
}

// tell viper to read in appConfig file and parse it
func readAndLoadUserConfig() {
  if err := userConfigViperInst.ReadInConfig(); err != nil {
    mlog.GetLogger().Error("ERROR READING USER CONF: ", err)

    switch err.(type) {
    case viper.ConfigFileNotFoundError:
      // pass if appConfig file does not exist
      mlog.GetLogger().Error("ERROR IS TYPED AS CONFIG FILE NOT FOUND ")
      err := createDefaultUserConfigFile()
      if err != nil {
        mlog.GetLogger().Error(fmt.Sprintf("Error creating default user config file, %s\n", err.Error()))
        os.Exit(1)
      }
    default:
      mlog.GetLogger().Error(fmt.Sprintf("Error reading user config file, %s\n", err.Error()))
      os.Exit(1)
    }
  }
  // parse appConfig file into their corresponding struct
  LoadUserConfig()
}

func createDefaultUserConfigFile() error {
  userConfigFilename := filepath.Join(CONFIG_PATH, "user_config."+CONFIG_EXT)
  f, err := os.Create(userConfigFilename)
  if err != nil {
    return err
  }
  f.Close()

  userConfigViperInst.Set("blockchain.networkContractAddress", mnet_vars.EMPTY_CONTRACT_ADDRESS)
  userConfigViperInst.Set("blockchain.chainId", 161027)

  return userConfigViperInst.WriteConfig()
}

// set the default value for the required configs, defaults are overridden by appConfig file if they exist
func setAppConfigDefaults() {
  // a list of seed nodes used for the initial peer discovery in the DHT (comma separated)
  appConfigViperInst.SetDefault("dht.bootnodes",
      "<NOT_CONFIGURED>:24801"
  )
  // max number of packets handled per second, increase this if bootnode is overwhelmed (set to negative for unlimited)
  appConfigViperInst.SetDefault("dht.maxIncomingPacketsPerSecond", 100000)
  // ignore a client's request packet if exceeds this limit, guard against spammy clients
  appConfigViperInst.SetDefault("dht.maxPerClientIncomingPacketsPerMinute", 500)
  // announce itself to the DHT every x seconds
  appConfigViperInst.SetDefault("dht.announceSelfIntervalSeconds", 5)
  // send find peers request to the DHT every x seconds
  appConfigViperInst.SetDefault("dht.requestPeersIntervalSeconds", 5)
  // determines whether DHT persist routing table to disk periodically and read it on startup
  appConfigViperInst.SetDefault("dht.cacheRoutingTableToDisk", false)

  appConfigViperInst.SetDefault("log.dir", ".")
  appConfigViperInst.SetDefault("log.level", "warn")

  appConfigViperInst.SetDefault("global.vlanFilterEnabled", false)
}

func setUserConfigDefaults() {
  userConfigViperInst.SetDefault("blockchain.networkContractAddress", mnet_vars.EMPTY_CONTRACT_ADDRESS)
  userConfigViperInst.SetDefault("blockchain.chainId", 161027)
}
