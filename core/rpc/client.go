package rpc

import (
  "bytes"
  "github.com/gorilla/rpc/v2/json2"
  "net/http"
)

type Client struct {
  URL string
}

func NewRPCClient(remoteHost string) *Client {
  return &Client{
    URL: "http://" + remoteHost + ":" + RPC_PORT + RPC_PATH,
  }
}

func (c Client) Call(method string, args interface{}, reply interface{}) error {
  message, err := json2.EncodeClientRequest(method, args)
  if err != nil {
    return err
  }

  resp, err := http.Post(c.URL, "application/json", bytes.NewReader(message))
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  return json2.DecodeClientResponse(resp.Body, &reply)
}
