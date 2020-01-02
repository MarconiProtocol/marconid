package ranker_strategy

import (
  "../../../../../core/runtime"
  "../utils"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "math"
  "net"
  "sort"
  "sync"
  "time"
)

const LATENCY_HISTORY_LEN = 50
const DEFAULT_DELAY_BETWEEN_LATENCY_CHECK_ICMP = 2000
const HEALTH_PORT = 23800

/*
	Ranking keeps track of node health statistics
	and helps to rank nodes based on those statistics
    implements Ranker interface
*/
type LatencyRankingStrategy struct {
  nodeIps                   map[string]bool
  nodeHealthStatistics      map[string]*nodeHealthStatistics
  healthChecksSignal        chan bool
  latencyChecksOnce         sync.Once
  nodeHealthStatisticsMutex sync.Mutex
}

/*
  Create ranking strategy
*/
func (l *LatencyRankingStrategy) Initialize() {
  l.nodeIps = make(map[string]bool)
  l.nodeHealthStatistics = make(map[string]*nodeHealthStatistics)
  l.nodeHealthStatisticsMutex = sync.Mutex{}
  l.StartLatencyChecks()
}

/*
	Add nodes to be tested by ranker
*/
func (l *LatencyRankingStrategy) AddNode(ipAddr string) {
  if _, exists := l.nodeIps[ipAddr]; !exists {
    l.nodeIps[ipAddr] = true
    mlog.GetLogger().Info(fmt.Sprintf("Added peer %s to latency ranker", ipAddr))
  }
}

/*
	remove nodes that should no longer be tested by ranker
*/
func (l *LatencyRankingStrategy) RemoveNode(ipAddr string) {
  if _, exists := l.nodeIps[ipAddr]; exists {
    delete(l.nodeIps, ipAddr)
    mlog.GetLogger().Info(fmt.Sprintf("Removed peer %s from latency ranker", ipAddr))
  }
}

/*
	Return the top N nodes
*/
func (l *LatencyRankingStrategy) GetTopNodes(num int) []string {
  // Get an array of node health stats, ordered by latency averages, smallest first
  sortedHealthStatistics := l.getHealthStatsSortedByLatencyAverages(true)
  numHealthStats := len(sortedHealthStatistics)

  // basically floor num to number of available nodes
  if num > numHealthStats {
    num = numHealthStats
  }

  // Return an array of nodes' IPs in latency average sorted order
  var results []string
  for i := 0; i < num; i++ {
    results = append(results, sortedHealthStatistics[i].IP)
  }
  return results
}

/*
	Performs a health check on nodes on an interval
*/
func (l *LatencyRankingStrategy) StartLatencyChecks() {
  l.latencyChecksOnce.Do(func() {
    go func() {
      mlog.GetLogger().Info("Starting health checks")

      for {
        select {
        case <-l.healthChecksSignal:
          mlog.GetLogger().Info("Terminating health checks")
          break
        default:
          if len(l.nodeIps) == 0 {
            mlog.GetLogger().Info("  Ranking::StartLatencyChecks() -> no nodes to rank yet")
          } else {
            mlog.GetLogger().Debug(fmt.Sprintf("Ranking::StartLatencyChecks()"))
            // do health check for nodes
            for nodeIp := range l.nodeIps {
              mlog.GetLogger().Debug(fmt.Sprintf("  Ranking::StartLatencyChecks() ranking node %s", nodeIp))
              l.checkLatencyWithTCPConn(nodeIp, HEALTH_PORT)
            }
          }
          time.Sleep(time.Second * time.Duration(2))
        }
      }
    }()
  })
}

/*
	Returns an array of nodeHealthStatistics sorted by latency averages
*/
func (l *LatencyRankingStrategy) getHealthStatsSortedByLatencyAverages(filterZeroAvgs bool) []*nodeHealthStatistics {
  // Create array of nodeHealthStatistics, possibly filtering zero or negative latency averages
  var nodeHealthStats []*nodeHealthStatistics
  l.nodeHealthStatisticsMutex.Lock()
  for _, phs := range l.nodeHealthStatistics {
    if !filterZeroAvgs || (filterZeroAvgs && phs.LatencyAverage != 0) {
      if isValidRemoteNode(phs.IP) {
        nodeHealthStats = append(nodeHealthStats, phs)
      }
    }
  }
  l.nodeHealthStatisticsMutex.Unlock()

  // Sort array by the latency average values, lowest first
  sort.SliceStable(nodeHealthStats, func(i int, j int) bool {
    return nodeHealthStats[i].LatencyAverage < nodeHealthStats[j].LatencyAverage
  })

  return nodeHealthStats
}

/*
	Use TCP connection to make an estimate on latency
*/
func (l *LatencyRankingStrategy) checkLatencyWithTCPConn(ipAddr string, port uint16) {
  localAddr := mruntime.GetMRuntime().InterfaceInfo.GetLocalMainInterfaceIpAddr()

  addrs, err := net.LookupHost(ipAddr)
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("checkLatencyWith : Error resolving %s. %s\n", ipAddr, err))
  }
  remoteAddr := addrs[0]

  latency := magent_edge_ranker_utils.TimeTCPSynAck(localAddr, remoteAddr, port)
  mlog.GetLogger().Debug(fmt.Sprintf("    Ranking::checkLatencyWithTCPConn node %s:%d latency %s", ipAddr, port, latency))
  l.updateNodeHealth(ipAddr, latency)
}

/*
	Use ICMP pings to make an estimate on latency
*/
func (l *LatencyRankingStrategy) checkLatencyWithPinger(ipAddr string) error {
  const NUM_PINGS = 3
  const TIMEOUT = time.Duration(time.Second * 1)

  pinger, err := magent_edge_ranker_utils.NewPinger(ipAddr)
  if err != nil {
    mlog.GetLogger().Fatal(fmt.Sprintf("rank.go::checkLatency failed %s", err.Error()))
    return fmt.Errorf("rank.go::checkLatency failed %s", err.Error())
  }
  pinger.Count = NUM_PINGS
  pinger.Timeout = TIMEOUT
  pinger.Run()

  stats := pinger.Statistics()

  // TODO double check that number of packets received is the same as configured?
  mlog.GetLogger().Info(fmt.Sprintf("checkLatencyWithPinger against %s, average RTT: %s", ipAddr, stats.AvgRtt.String()))

  l.updateNodeHealth(ipAddr, stats.AvgRtt)
  return nil
}

/*
	Update the nodes's health statistics with a new latency entry
*/
func (l *LatencyRankingStrategy) updateNodeHealth(ipAddr string, latency time.Duration) {
  l.nodeHealthStatisticsMutex.Lock()
  if _, exists := l.nodeHealthStatistics[ipAddr]; !exists {
    l.nodeHealthStatistics[ipAddr] = newPeerHealthStatistic(ipAddr)
  }
  l.nodeHealthStatisticsMutex.Unlock()
  nodeHealth := l.nodeHealthStatistics[ipAddr]
  nodeHealth.addLatency(latency)
}

func isValidRemoteNode(remoteHost string) bool {
  if "" == remoteHost ||
    "127.0.0.1" == remoteHost {
    return false
  }
  return true
}

/*
	Data object storing a node's health statistics,
	keeping track of LATENCY_HISTORY_LEN number of latency entries and a latency mean
*/
type nodeHealthStatistics struct {
  IP             string
  LatencyHistory []int64 // A rolling window of latency entries,
  // The index resets to 0 as it reaches LATENCY_HISTORY_LEN, allow us to overwrite old entries.
  LatencyAverage int64 // Simple average (mean) of the latency entries ms
}

func newPeerHealthStatistic(peerIp string) *nodeHealthStatistics {
  peerHealth := nodeHealthStatistics{
    peerIp,
    []int64{},
    0,
  }
  return &peerHealth
}

/*
	Add a latency entry to a node's health statistic record, as well as updating the latency average
*/
func (p *nodeHealthStatistics) addLatency(latency time.Duration) {
  latencyNanoseconds := float64(latency.Nanoseconds())
  newAverage := float64(p.LatencyAverage)
  if len(p.LatencyHistory) >= LATENCY_HISTORY_LEN {
    newAverage -= float64(p.LatencyHistory[0]) / float64(LATENCY_HISTORY_LEN)
    p.LatencyHistory = p.LatencyHistory[1:]
  } else {
    newAverage *= float64(len(p.LatencyHistory)) / float64(len(p.LatencyHistory)+1)
  }
  p.LatencyHistory = append(p.LatencyHistory, latency.Nanoseconds())
  newAverage += latencyNanoseconds / float64(len(p.LatencyHistory))
  p.LatencyAverage = int64(math.Max(math.Round(newAverage), 1))
  mlog.GetLogger().Debug(fmt.Sprintf("    Node %s new average latency %d", p.IP, p.LatencyAverage))
}
