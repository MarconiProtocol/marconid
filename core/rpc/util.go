package rpc

import (
  "errors"
  "fmt"
  "strings"
)

func ParseIPAndPort(host string) (string, string, error) {
  parsedStr := strings.Split(host, ":")
  if len(parsedStr) != 2 {
    return "", "", errors.New(fmt.Sprintf("Could not parse the string: %s into an ip and port", host))
  }
  return parsedStr[0], parsedStr[1], nil
}
