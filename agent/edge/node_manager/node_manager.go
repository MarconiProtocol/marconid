package magent_edge_node_manager

import (
  "../../../core/rpc"
  "../../service/base"
  "./ranker"
  "./ranker/ranker_strategy"
  "errors"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "sort"
  "sync"
  "time"
)

/*
  The edge node finder keeps track of nodes that have been discovered through the edge route beacon,
  who have successfully passed a health check.
*/
type EdgeNodeManager struct {
  nodes   map[string]*EdgeNode // IP to node
  rankers []*ranker.Ranker

  nodesMutex sync.Mutex
}

type EdgeNode struct {
  Ip                  string
  PubKeyHash          string
  HealthCheckFailures int
  LastChecked         int64
}

const MAX_FAILED_HEALTH_CHECKS = 5
const BAN_DURATION_S = 300
const NODE_HEALTH_CHECK_INTERVAL_S = 5

var instance *EdgeNodeManager
var once sync.Once

func Instance() *EdgeNodeManager {
  once.Do(func() {
    instance = &EdgeNodeManager{}
    instance.initialize()
  })
  return instance
}

func (e *EdgeNodeManager) initialize() {
  e.nodes = make(map[string]*EdgeNode)

  // TODO tyler - make this dynamic when UI is added.
  e.rankers = append(e.rankers, ranker.NewRanker(&ranker_strategy.LatencyRankingStrategy{}))

  // Start health checks against nodes that have been added to the finder
  e.StartHealthChecks()
}

/*
  Add a node to the node manager. This creates an EdgeNode instance to the nodes dictionary.
  These nodes are processed on an interval through the goroutine periodicallyPerformHealthChecksOnAllNodes
*/
func (e *EdgeNodeManager) AddNode(ip string, expectedPkh string) {
  e.nodesMutex.Lock()
  // no-op if the node has already been added
  if _, exists := e.nodes[ip]; !exists {
    node := EdgeNode{
      ip,
      expectedPkh,
      0,
      0,
    }
    e.nodes[ip] = &node
  }
  e.nodesMutex.Unlock()
}

/*
  Obtain a slice of "top" nodes from the ranker.
  The top heuristic SHOULD be based on whatever rankers are provided as the strategy - but this definition belongs to the ranker
*/
func (e *EdgeNodeManager) GetTopNodes(numNodes int) []string {
  nodeScore := map[string]int{}
  for _, rank := range e.rankers {
    ordering := rank.GetTopNodes(numNodes)
    // Add score to ranked nodes
    for score, node := range ordering {
      nodeScore[node] += score
    }
    // Add worst score to unranked nodes
    for node, _ := range nodeScore {
      found := false
      for _, ip := range ordering { // Contains
        if node == ip {
          found = true
          break
        }
      }
      if !found {
        nodeScore[node] += numNodes + 1
      }
    }
  }
  sortedNodes := e.sortNodesByScore(nodeScore)
  return sortedNodes
}

/*
  Sort map keys by the value
*/
func (e *EdgeNodeManager) sortNodesByScore(nodeScore map[string]int) []string {
  scoreToIpMap := map[int][]string{}
  var scoreList []int
  for nodeIp, score := range nodeScore {
    scoreToIpMap[score] = append(scoreToIpMap[score], nodeIp)
  }
  for score := range scoreToIpMap {
    scoreList = append(scoreList, score)
  }
  sort.Sort(sort.Reverse(sort.IntSlice(scoreList)))
  var sortedNodeIps []string
  for _, scoring := range scoreList {
    for _, nodeIp := range scoreToIpMap[scoring] {
      sortedNodeIps = append(sortedNodeIps, nodeIp)
    }
  }
  return sortedNodeIps
}

/*
  Waits until there are 'numNodes' number of top nodes returned by the ranker.
  Retests at an interval of searchIntervalMS
*/
func (e *EdgeNodeManager) WaitForTopNodes(numNodes int, searchIntervalMS int) []string {
  var topPeers []string
  for {
    topPeers = e.GetTopNodes(numNodes)
    // wait for more than 1 peer, for now
    if len(topPeers) >= numNodes {
      break
    }
    mlog.GetLogger().Info(fmt.Sprintf("EdgeNodeManager::WaitForTopNodes - Not enough nodes returned %d/%d, retrying.", len(topPeers), numNodes))
    time.Sleep(time.Millisecond * time.Duration(searchIntervalMS))
  }
  return topPeers
}

/*
  Starts a goroutine that will call performHealthCheck on an interval, adding and removing nodes as necessary from the ranker.
*/
func (e *EdgeNodeManager) StartHealthChecks() {
  go func() {
    for {
      // Iterate through all nodes, and perform health checks
      e.nodesMutex.Lock()
      for nodeIp, node := range e.nodes {
        err := e.performHealthCheck(node)
        if err != nil {
          // If the node has failed too many times, we will remove it from the ranker
          if node.HealthCheckFailures >= MAX_FAILED_HEALTH_CHECKS {
            e.removeFromRankers(nodeIp)
          }
        } else {
          // If there is no error, we can attempt to add it to the ranker
          // if the ranker already knows about this node, it is a no-op
          e.addToRankers(nodeIp)
        }

      }
      e.nodesMutex.Unlock()
      time.Sleep(time.Second * time.Duration(NODE_HEALTH_CHECK_INTERVAL_S))
    }
  }()
}

func (e *EdgeNodeManager) addToRankers(ipAddr string) {
  for _, ranker := range e.rankers {
    ranker.AddNode(ipAddr)
  }
}

func (e *EdgeNodeManager) removeFromRankers(ipAddr string) {
  for _, ranker := range e.rankers {
    ranker.RemoveNode(ipAddr)
  }
}

/*
  Performs a health check against an edge node.
  If the health check passes, no error is returned, otherwise an error is returned
*/
func (e *EdgeNodeManager) performHealthCheck(edgeNode *EdgeNode) error {

  // We want to retest this edge node as a potentially good node under the following conditions
  // We need to consider the following two cases:
  //    1. If the edgeNode has not yet exceeded the max checks
  //    2. If it has already served its ban
  // If it meets the above conditions, we should see if the node is a good node
  // Otherwise, the edge node is not eligible for a check due to it being banned or above max checks, so we just return an error
  now := time.Now().Unix()
  var nodeCheckError error = nil
  if edgeNode.HealthCheckFailures < MAX_FAILED_HEALTH_CHECKS || edgeNode.LastChecked+BAN_DURATION_S < now {
    // Send a health check RPC
    edgeNode.LastChecked = now
    pkh, err := e.sendNodeFinderCheckRPC(edgeNode.Ip)
    if err != nil {
      edgeNode.HealthCheckFailures += 1
      nodeCheckError = errors.New(fmt.Sprintf("Health check failed for node %s. Checks: %d. Err: %s", edgeNode.Ip, edgeNode.HealthCheckFailures, err))
    } else if pkh != edgeNode.PubKeyHash {
      // this disqualifies the node, because it responded and told us it is broadcasting for another network
      edgeNode.HealthCheckFailures = MAX_FAILED_HEALTH_CHECKS
      nodeCheckError = errors.New(fmt.Sprintf("Health check failed for node %s. It is not broadcasting the pkh we expected. DHT cache could be out of date. Expected: %s, Actual: %s", edgeNode.Ip, edgeNode.PubKeyHash, pkh))
    } else {
      // good case, we want to reset its checks
      edgeNode.HealthCheckFailures = 0
    }
  } else {
    nodeCheckError = errors.New(fmt.Sprintf("Skipped health check for node %s. Checks: %d. Banned until: %d", edgeNode.Ip, edgeNode.HealthCheckFailures, edgeNode.LastChecked+BAN_DURATION_S))
  }
  return nodeCheckError
}

/*
  Performs an RPC to check if the remote node is available. It will provide the pkh of its edge route beacon key on return.
*/
func (e *EdgeNodeManager) sendNodeFinderCheckRPC(target string) (string, error) {
  reply := magent_base.NodeFinderCheckReply{}
  client := rpc.NewRPCClient(target)
  if err := client.Call("AgentClientService.NodeFinderCheckRPC", nil, &reply); err != nil {
    return "", errors.New(fmt.Sprintf("sendNodeFinderCheckRPC - Nil/error response from: %s; response: %v", target, err))
  }
  return reply.EdgeBeaconKey, nil
}
