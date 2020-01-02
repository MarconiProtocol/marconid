package mpacket_filter

import (
  "../../../../util"
  "../../../config"
  "fmt"
  "github.com/MarconiProtocol/gopacket"
  mlog "github.com/MarconiProtocol/log"
  m_packet_filter "github.com/MarconiProtocol/sdk/packet/filter"
  "io/ioutil"
  "path/filepath"
  "plugin"
  "sync"
)

/*
  Stores the result of filtering packet and the filters response that caused it to be dropped
*/
type PacketFilterResult struct {
  Action         m_packet_filter.Action
  FilterResponse *m_packet_filter.FilterResponse
}

/*
  The packet filter manager
*/
type PacketFilterManager struct {
  availableFilters map[string]string
  usedFilters      map[string]bool
  loadedFilters    map[string]*m_packet_filter.Filter
}

var filterManager *PacketFilterManager
var once sync.Once

/*
  Singleton getter
*/
func GetFilterManagerInstance() *PacketFilterManager {
  once.Do(func() {
    filterManager = &PacketFilterManager{}
    filterManager.init()
    filterManager.loadConfig()
    filterManager.loadFilter()
  })
  return filterManager
}

/*
  Initializes the properties used by the PacketFilterManager
*/
func (filterManager *PacketFilterManager) init() {
  mlog.GetLogger().Debug("packetFilterManager::init")
  filterManager.availableFilters = map[string]string{}
  filterManager.usedFilters = make(map[string]bool, 0)
  filterManager.loadedFilters = map[string]*m_packet_filter.Filter{}
}

/*
  Loads config.yml and based on the packet filter directory configured,
  reads the directory and caches a map of potential filter plugins to load
*/
func (filterManager *PacketFilterManager) loadConfig() {
  files, err := ioutil.ReadDir(mconfig.GetAppConfig().Global.Packet_Filter_Data_Directory_Path)
  if err != nil {
    mlog.GetLogger().Error(fmt.Sprintf("packetFilterManager::loadConfig - cannot read packet filter directory %v", err))
    mlog.GetLogger().Error("packetFilterManager::loadConfig - ", mconfig.GetAppConfig().Global.Packet_Filter_Data_Directory_Path)
  } else {
    for _, file := range files {
      if !file.IsDir() {
        filterManager.availableFilters[file.Name()] = file.Name()
      }
    }
  }
}

/*
  Attempts to load all potential plugins cached in availableFilters and cache the loaded filter object
*/
func (filterManager *PacketFilterManager) loadFilter() {
  mlog.GetLogger().Debug("packetFilterManager::load")
  for filterName, soFilePath := range filterManager.availableFilters {

    filterPath := filepath.Join(mconfig.GetAppConfig().Global.Packet_Filter_Data_Directory_Path, soFilePath)
    mlog.GetLogger().Debug(fmt.Sprintf("packetFilterManager::load - loading %s at %s", filterName, filterPath))

    if mutil.DoesExist(filterPath) {
      if filterPlugin, err := plugin.Open(filterPath); err == nil {
        // Lookup the object named mnet_filters.FILTER_INTERFACE_OBJ_NAME
        filterInterfaceSymbol, err := filterPlugin.Lookup(m_packet_filter.FILTER_INTERFACE_OBJ_NAME)
        if err != nil {
          mlog.GetLogger().Fatal(fmt.Sprintf("Failed to find symbol %s in the loaded plugin %s: %v", m_packet_filter.FILTER_INTERFACE_OBJ_NAME, filterName, err))
        }
        // Type assert the symbol to the interface type mnet_filters.Filter
        filterInterface, ok := filterInterfaceSymbol.(m_packet_filter.Filter)
        if !ok {
          mlog.GetLogger().Fatal(fmt.Sprintf("Failed to assert %s symbol to Filter interface type", m_packet_filter.FILTER_INTERFACE_OBJ_NAME))
        }
        // Store the loaded Filter interface object
        filterManager.loadedFilters[filterName] = &filterInterface
        mlog.GetLogger().Debug("packetFilterManager::load - loaded : " + filterName)

      } else {
        mlog.GetLogger().Error("packetFilterManager::load - error loading : "+filterName, err)
      }
    } else {
      mlog.GetLogger().Error(fmt.Sprintf("packetFilterManager::load - error loading : %s, filter not found at path %s ", filterName, filterPath))
    }
  }
}

/*
  Applies all loaded Filters against the provided packet
  When all filters have completed the result will be pushed to signalChannel
  This is a blocking call

  Filters are run concurrently using goroutines and a WaitGroup is used to wait for results from all filters

  Try to do as little work as possible in this function
*/
func (filterManager *PacketFilterManager) ProcessPacket(packet *gopacket.Packet, resultChannel *chan PacketFilterResult) {
  numFilters := len(filterManager.loadedFilters)

  wg := sync.WaitGroup{}
  wg.Add(numFilters)

  // TODO: may need to create a pool of filterResponse that we can just reuse,
  // but there is a risk people do not write to Msg or any future properties, (stale data)
  responses := make([]m_packet_filter.FilterResponse, numFilters)

  // apply each filter to the packet
  i := 0
  for _, filter := range filterManager.loadedFilters {
    go runFilter(*packet, filter, &responses[i], &wg)
    i++
  }

  // wait until all filters are completed
  wg.Wait()

  // Note: to avoid passing complexity to the end user, the current assumption is that there is no need to configure
  // the loaded plugins.
  //
  // We want to avoid the need for the end user having to define priorities for filters and/or logical operations applied
  // to the results from multiple filters
  //
  // KISS strategy for now where we drop packets if ANY filter returns a DROP signal
  result := calculateFinalResult(responses...)
  *resultChannel <- result
}

/*
  Return the final action to take given the responses from a collection of FilterResponses
*/
func calculateFinalResult(responses ...m_packet_filter.FilterResponse) PacketFilterResult {
  // Iterate through all responses, if we find a DROP action, return DROP early as well as the FilterResponse
  for _, response := range responses {
    if response.Status == m_packet_filter.DROP {
      return PacketFilterResult{m_packet_filter.DROP, &response}
    }
  }
  // If none of the responses were a DROP, then we return an ACCEPT
  return PacketFilterResult{m_packet_filter.ACCEPT, nil}
}

/*
  Invoke the Execute function of a Filter and decrement the WaitGroup to signify completion
*/
func runFilter(packet gopacket.Packet, filter *m_packet_filter.Filter, response *m_packet_filter.FilterResponse, wg *sync.WaitGroup) {
  defer wg.Done()
  (*filter).Execute(packet, response)
}
