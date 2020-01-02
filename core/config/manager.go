package mconfig

import (
  "../net/vars"
  "../rpc"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "github.com/fsnotify/fsnotify"
  "github.com/pkg/errors"
  "github.com/spf13/viper"
  "net/http"
  "os"
  "path/filepath"
)

const CONFIG_PATH = "/etc/marconid"
const CONFIG_EXT = "yml"

var appConfigViperInst *viper.Viper
var userConfigViperInst *viper.Viper

func InitializeConfigs(baseDir string) {
  InitializeAppConfig(baseDir)
  InitializeUserConfig(baseDir)
  userConfigService := new(UserConfigService)
  if err := rpc.RegisterService(userConfigService, ""); err != nil {
    mlog.GetLogger().Error("Failed to register RPC service for UserConfig")
  }
}

// initialize the appConfig object by parsing the appConfig file each time it's updated
// config_manager.go is kept outside of package "appConfig" on purpose so that we can avoid cyclic reference
// for packages that need to access the appConfig object, yet also expect to be notified when appConfig file changes
func InitializeAppConfig(baseDir string) {
  appConfigViperInst = viper.New()

  configName := "config"
  appConfigViperInst.SetConfigName(configName)
  appConfigViperInst.SetConfigType(CONFIG_EXT)
  appConfigViperInst.AddConfigPath(baseDir + CONFIG_PATH)

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

func InitializeUserConfig(baseDir string) {
  userConfigViperInst = viper.New()

  configName := "user_config"
  userConfigViperInst.SetConfigName(configName)
  userConfigViperInst.SetConfigType(CONFIG_EXT)
  userConfigViperInst.AddConfigPath(baseDir + CONFIG_PATH)

  // set default value for required configs
  setUserConfigDefaults(baseDir)

  // read in user_config.yml
  readAndLoadUserConfig(baseDir)

  // triggered on initial read and every time user_config.yml is modified
  // make sure any functions called within OnConfigChange are thread-safe
  userConfigViperInst.OnConfigChange(func(event fsnotify.Event) {
    // read in user_config.yml
    readAndLoadUserConfig(baseDir)
    fmt.Println("Config file changed:", event.Name)
  })
  // watch the user_config.yml for file change
  userConfigViperInst.WatchConfig()
}

// tell viper to read user_config.yml and parse it
func readAndLoadUserConfig(baseDir string) {
  if err := userConfigViperInst.ReadInConfig(); err != nil {
    switch err.(type) {
    case viper.ConfigFileNotFoundError:
      err := createDefaultUserConfigFile(baseDir)
      if err != nil {
        mlog.GetLogger().Error(fmt.Sprintf("Error creating default user config file, %s\n", err.Error()))
        os.Exit(1)
      }
    default:
      mlog.GetLogger().Error(fmt.Sprintf("Error reading user config file, %s\n", err.Error()))
      os.Exit(1)
    }
  }
  // parse user_config.yml into its corresponding struct
  LoadUserConfig()
}

func createDefaultUserConfigFile(baseDir string) error {
  userConfigFilename := filepath.Join(baseDir, CONFIG_PATH, "user_config."+CONFIG_EXT)
  setUserConfigDefaults(baseDir)
  return userConfigViperInst.WriteConfigAs(userConfigFilename)
}

// set the default value for the required configs, defaults are overridden by appConfig file if they exist
func setAppConfigDefaults() {
  // max number of packets handled per second, increase this if bootnode is overwhelmed (set to negative for unlimited)
  appConfigViperInst.SetDefault("dht.max_incoming_packets_per_second", 100000)
  // ignore a client's request packet if exceeds this limit, guard against spammy clients
  appConfigViperInst.SetDefault("dht.max_per_client_incoming_packets_per_minute", 500)
  // announce itself to the DHT every x seconds
  appConfigViperInst.SetDefault("dht.announce_base_interval_seconds", 60)
  // announce itself to the DHT every x seconds
  appConfigViperInst.SetDefault("dht.announce_self_interval_seconds", 15)
  // send find peers request to the DHT every x seconds
  appConfigViperInst.SetDefault("dht.request_peers_interval_seconds", 5)
  // determines whether DHT persist routing table to disk periodically and read it on startup
  appConfigViperInst.SetDefault("dht.cache_routing_table_to_disk", true)

  appConfigViperInst.SetDefault("log.dir", ".")
  appConfigViperInst.SetDefault("log.level", "warn")

  appConfigViperInst.SetDefault("global.vlan_filter_enabled", false)

  appConfigViperInst.SetDefault("global.bridge_ageing_time_seconds", 30)
}

func setUserConfigDefaults(baseDir string) {
  userConfigViperInst.SetDefault("global.base_dir", baseDir)
  userConfigViperInst.SetDefault("blockchain.network_contract_address", mnet_vars.EMPTY_CONTRACT_ADDRESS)
}

type UserConfigService struct{}

type UpdateNetworkContractAddressArgs struct {
  NetworkContractAddress string
}

type UpdateNetworkContractAddressReply struct{}

func (u *UserConfigService) UpdateNetworkContractAddressRPC(r *http.Request, args *UpdateNetworkContractAddressArgs, reply *UpdateNetworkContractAddressReply) error {
  // if the folder does not exist
  if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
    if err1 := os.MkdirAll(CONFIG_PATH, os.ModePerm); err1 != nil {
      return err1
    }
  }

  // write to user_conf.yml
  if GetUserConfig().Blockchain.Network_Contract_Address == mnet_vars.EMPTY_CONTRACT_ADDRESS {
    userConfigViperInst.Set("blockchain.network_contract_address", args.NetworkContractAddress)
    if err := userConfigViperInst.WriteConfig(); err != nil {
      return err
    } else {
      return nil
    }
  } else {
    return errors.New("Join network skipped, user is already part of network: " + GetUserConfig().Blockchain.Network_Contract_Address)
  }
}
