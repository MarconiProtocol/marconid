package magent_middleware_interface

import (
  "../../../../core/config"
  "../../../../core/rpc"
  "bytes"
  "encoding/json"
  "gitlab.neji.vm.tc/marconi/log"
  "io/ioutil"
  "net/http"
)

const (
  HOST        = "http://127.0.0.1"
  PORT        = "28902"
  MARCONID_ID = "18"
)

type SubscribeParams struct {
  Id                     string
  Port                   string
  NetworkContractAddress string
  PubKeyHash             string
}

type UnsubscribeParams struct {
  Id string
}

/*
  Marshal SubscribeParams into JSON string format
*/
func (r SubscribeParams) GetParamsJsonString() string {
  paramsBytes, err := json.Marshal(r)
  if err != nil {
    mlog.GetLogger().Fatal(err)
  }
  return string(paramsBytes)
}

/*
  Marshal UnsubscribeParams into JSON string format
*/
func (r UnsubscribeParams) GetParamsJsonString() string {
  paramsBytes, err := json.Marshal(r)
  if err != nil {
    mlog.GetLogger().Fatal(err)
  }
  return string(paramsBytes)
}

/*
  Create a JSON RPC payload given a method name and params for that RPC
*/
func createJsonPayload(method string, params string) []byte {
  payloadStr := "{\"jsonrpc\":\"2.0\", \"id\":1, \"method\":\"" + method + "\", \"params\":" + params + " }"
  return []byte(payloadStr)
}

/*
  Register to middleware for peer updates
*/
func RegisterForPeerUpdates(pubKeyHash string) error {
  // TODO: ayuen, we need to read / get  the network cluster key from somewhere
  networkContractAddress := mconfig.GetUserConfig().Blockchain.NetworkContractAddress
  path := "api/middleware/v1"

  subscribeParams := SubscribeParams{
    MARCONID_ID,
    rpc.RPC_PORT,
    networkContractAddress,
    pubKeyHash,
  }
  paramsStr := subscribeParams.GetParamsJsonString()
  payload := createJsonPayload("subscribe", paramsStr)
  mlog.GetLogger().Debug(string(payload))
  err := sendJsonRpcOverHttp(HOST, PORT, path, payload)
  return err
}

/*
  Unregister from middleware for peer updates
*/
func UnregisterForPeerUpdates() error {
  path := "api/middleware/v1"

  unsubscribeParams := UnsubscribeParams{
    MARCONID_ID,
  }
  paramsStr := unsubscribeParams.GetParamsJsonString()
  payload := createJsonPayload("unsubscribe", paramsStr)
  mlog.GetLogger().Info(string(payload))
  err := sendJsonRpcOverHttp(HOST, PORT, path, payload)
  return err
}

/*
  Send JSON RPC over http to a middleware process
*/
func sendJsonRpcOverHttp(host string, port string, path string, jsonPayloadBytes []byte) error {
  url := host + ":" + port + "/" + path
  request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayloadBytes))
  request.Header.Set("Content-Type", "application/json")

  client := new(http.Client)
  response, err := client.Do(request)
  if err != nil {
    mlog.GetLogger().Error(err)
    return err
  }
  defer response.Body.Close()

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    mlog.GetLogger().Error(err)
    return err
  }
  mlog.GetLogger().Info(string(body))

  return err
}
