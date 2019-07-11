package rpc

import (
  "../../util"
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/log"
  "io/ioutil"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "sync"
)

const (
  REQUEST_EDGE_ROUTE_PORT          = "rpcRequestEdgeRoutePort"
  REQUEST_INTERNAL_EDGE_ROUTE_PORT = "rpcRequestInternalEdgeRoutePort"
  REQUEST_DH_KEY_EXCHANGE_SYN      = "rpcDHKeyExchangeSYN"
  REQUEST_DH_KEY_EXCHANGE_ACK      = "rpcDHKeyExchangeACK"
  UPDATE_PEERS                     = "rpcUpdatePeers"
  UPDATE_EDGE_PEERS                = "rpcUpdateEdgePeers"
  REQUEST_PUB_KEY_EXCHANGE_SYN     = "rpcPubKeyExchangeSYN"
  REQUEST_PUB_KEY_EXCHANGE_ACK     = "rpcPubKeyExchangeACK"
)

const (
  RPC_PORT = "24802"
)

var rpcHandlerMapMutex sync.Mutex
var rpcHandlerMap = make(map[string]func(*http.Request, http.ResponseWriter, string, string))

func StartRPCServer() {
  StartRPCServerOnPort(24802)
}

func StartRPCServerOnPort(port int) {
  fmt.Println("Starting HTTP RPC Server on port: ", port)

  http.HandleFunc("/rpc/m/request", RouteRPC)
  http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func RegisterRpcHandler(rpcName string, handler func(*http.Request, http.ResponseWriter, string, string)) {
  rpcHandlerMapMutex.Lock()
  if _, exists := rpcHandlerMap[rpcName]; !exists {
    rpcHandlerMap[rpcName] = handler
  }
  rpcHandlerMapMutex.Unlock()
}

func SendRPC(remoteHost string, remotePort string, rpcType string, rpcPayload string) *RPCResponse {
  var query bytes.Buffer
  query.WriteString("http://")
  query.WriteString(remoteHost)
  query.WriteString(":")
  query.WriteString(remotePort)
  query.WriteString("/rpc/m/request?m=")

  payload := buildRPCPayload(rpcType, rpcPayload)

  encodedPayload := mutil.EncodeBase64(payload)
  safeEncodedPayload := url.QueryEscape(string(encodedPayload))
  query.WriteString(safeEncodedPayload)

  resp, err := http.Get(query.String())
  if err != nil {
    msg := fmt.Sprintf("RPC send failure %s", err)
    mlog.GetLogger().Error(msg)
    return &RPCResponse{ Error: msg }
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  response := parseResponse(string(body))

  return response
}

func RouteRPC(w http.ResponseWriter, r *http.Request) {
  encodedPayload := r.FormValue("m")
  if encodedPayload != "" {

    encPayload, err := mutil.DecodeBase64([]byte(encodedPayload))
    if err != nil {
      mlog.GetLogger().Error("routeRPC: failed to decode: ", encodedPayload)
    }
    payload := encPayload

    data := strings.Split(string(payload[:]), "#M#")
    if len(data) > 2 {
      rpc := data[0]
      reqInfohash := data[1]
      reqPayload := data[2]
      callRPC(r, w, rpc, reqInfohash, reqPayload)
    }
  }
}

func callRPC(r *http.Request, w http.ResponseWriter, rpcName, reqInfohash, reqPayload string) {
  rpcHandlerMapMutex.Lock()
  handler, exists := rpcHandlerMap[rpcName]
  rpcHandlerMapMutex.Unlock()

  if exists {
    mlog.GetLogger().Debug(fmt.Sprintf("RPC::callRPC - Routing rpc to handler for %s", rpcName))
    handler(r, w, reqInfohash, reqPayload)
  } else {
    mlog.GetLogger().Warning(fmt.Sprintf("RPC::callRPC - Did not find an rpc handler for %s", rpcName))
  }
}

func buildRPCPayload(rpcType string, rpcPayload string) []byte {
  var payload bytes.Buffer
  payload.WriteString(rpcType)
  payload.WriteString("#M#")
  payload.WriteString("someInfoHash")
  payload.WriteString(rpcType)
  payload.WriteString("#M#")
  payload.WriteString(rpcPayload)
  return payload.Bytes()
}

/*
  Takes in a string in the format of IP:PORT and returns both values separately.
  If the input string could not be parsed, returns an Error
*/
func ParseIPAndPort(host string) (string, string, error) {
  parsedStr := strings.Split(host, ":")
  if len(parsedStr) != 2 {
    return "", "", errors.New(fmt.Sprintf("Could not parse the string: %s into an ip and port", host))
  }
  return parsedStr[0], parsedStr[1], nil
}

type RPCResponse struct {
  Result string
  Error  string
}

func WriteResponseWithResult(w *http.ResponseWriter, result string) error {
  responseObj := RPCResponse{
    Result: result,
    Error:  "",
  }

  responseBytes, err := json.Marshal(responseObj)
  _, err = (*w).Write([]byte(responseBytes))
  if err != nil {
    return err
  }

  return nil
}

func WriteResponseWithError(w *http.ResponseWriter, msg string) error {
  responseObj := RPCResponse{
    Result: "",
    Error:  msg,
  }

  responseBytes, err := json.Marshal(responseObj)
  _, err = (*w).Write([]byte(responseBytes))
  if err != nil {
    return err
  }

  return nil
}

func parseResponse(response string) *RPCResponse {
  var responseObj RPCResponse
  if response != "" {
    responseBytes := []byte(response)
    err := json.Unmarshal(responseBytes, &responseObj)
    if err != nil {
      mlog.GetLogger().Error(fmt.Sprintf("Failed to parse rpc response - %s", err))
      return nil
    }
  } else {
    responseObj.Error = "Empty response from target"
  }

  return &responseObj
}
