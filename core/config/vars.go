package mconfig


type AppConfiguration struct {
	Global     GlobalConfiguration
	Blockchain BlockchainConfiguration
	DHT        DHTConfiguration
	Log        LogConfiguration
}

type UserConfiguration struct {
	Blockchain BlockchainUserConfiguration
}

type GlobalConfiguration struct {
	PacketFilterDataDirectoryPath string // default data directory path for packet filter functions
	PacketFiltersEnabled          bool
	VlanFilterEnabled             bool // whether to enable vlan filter
}

type BlockchainConfiguration struct {
	StaticPeers            bool
}

type BlockchainUserConfiguration struct {
	ChainID                int
	NetworkContractAddress string
}

type DHTConfiguration struct {
	BootNodes                            string // a list of seed nodes used for the initial peer discovery in the DHT
	MaxIncomingPacketsPerSecond          int64  // max number of packets handled per second, increase this if bootnode is overwelmed (set to negative for unlimited)
	MaxPerClientIncomingPacketsPerMinute int    // ignore a client's request packet if exceeds this limit, guard against spammy clients
	AnnounceSelfIntervalSeconds          int    // announce itsef to the DHT every x seconds
	RequestPeersIntervalSeconds          int    // send find peers request to the DHT every x seconds
	CacheRoutingTableToDisk              bool   // determines whether DHT persist routing table to disk periodically and read it on startup
}

type LogConfiguration struct {
	Dir   string // directory which to store logs in
	Level string // dynamic (can be changed in appConfig.yml during runtime)
}