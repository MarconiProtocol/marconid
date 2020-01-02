package mconfig

type AppConfiguration struct {
  Global     GlobalConfiguration
  Blockchain BlockchainConfiguration
  DHT        DHTConfiguration
  Log        LogConfiguration
}

type UserConfiguration struct {
  Global     GlobalUserConfiguration
  Blockchain BlockchainUserConfiguration
}

type GlobalUserConfiguration struct {
  Base_Dir string // base directory for Marconid
}

type GlobalConfiguration struct {
  Packet_Filter_Data_Directory_Path string // default data directory path for packet filter functions
  Packet_Filters_Enabled            bool
  Vlan_Filter_Enabled               bool // whether to enable vlan filter
  BRIDGE_AGEING_TIME_SECONDS        int  // ageing time for marconi bridge
}

type BlockchainConfiguration struct {
  Chain_ID     int
  Static_Peers bool
}

type BlockchainUserConfiguration struct {
  Network_Contract_Address string
}

type DHTConfiguration struct {
  BootNodes                                  string // a list of seed nodes used for the initial peer discovery in the DHT
  Max_Incoming_Packets_Per_Second            int64  // max number of packets handled per second, increase this if bootnode is overwelmed (set to negative for unlimited)
  Max_Per_Client_Incoming_Packets_Per_Minute int    // ignore a client's request packet if exceeds this limit, guard against spammy clients
  Announce_Base_Interval_Seconds             int    // announce the base route
  Announce_Self_Interval_Seconds             int    // announce itsef to the DHT every x seconds
  Request_Peers_Interval_Seconds             int    // send find peers request to the DHT every x seconds
  Cache_Routing_Table_To_Disk                bool   // determines whether DHT persist routing table to disk periodically and read it on startup
}

type LogConfiguration struct {
  Dir   string // directory which to store logs in
  Level string // dynamic (can be changed in appConfig.yml during runtime)
}
