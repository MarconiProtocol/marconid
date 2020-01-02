package mutil

import (
  "../core/sys/cmd"
  "fmt"
  "reflect"
  "testing"
)

func TestExecuteCommand(t *testing.T) {
  t.Log("TestExecuteCommand called")
  result := msys_cmd.ExecuteCommand("ls", []string{"-al", "/tmp"})
  fmt.Println("TestExecuteCommand called: ", result)
}

func TestExecuteCommandPipe(t *testing.T) {
  t.Log("TestExecuteCommandPipe called")
  result := msys_cmd.ExecuteCommandPipe("ls", []string{"-al", "/tmp"}, "grep", []string{"go-build"})
  fmt.Println("TestExecuteCommandPipe called: ", string(result.Bytes()))
}

func TestExecuteCommandPipe_sucess1(t *testing.T) {
  t.Log("TestExecuteCommandPipe called")
  result := msys_cmd.ExecuteCommandPipe("ps", []string{"agx"}, "grep", []string{"bash"})
  fmt.Println("TestExecuteCommandPipe called: ", string(result.Bytes()))
}

func TestExecuteSequencialIdenticalCommand(t *testing.T) {
  t.Log("TestExecuteSequencialIdenticalCommand called")
  args := make(map[int][]string)
  args[0] = []string{"this is 1st command with arg1"}
  args[1] = []string{"this is 2nd command with arg1"}
  result := msys_cmd.ExecuteSequencialIdenticalCommand("echo", args)
  fmt.Println("TestExecuteSequencialIdenticalCommand called: ", result)
  for i, v := range result {
    fmt.Println("TestExecuteSequencialIdenticalCommand results: ", i, v)
  }
}

func TestGetMutualMPipePort(t *testing.T) {
  type args struct {
    pubKeyHash     string
    peerPubKeyHash string
  }
  tests := []struct {
    name string
    args args
    want int
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := GetMutualMPipePort(tt.args.pubKeyHash, tt.args.peerPubKeyHash); got != tt.want {
        t.Errorf("getMutualMPipePort() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestEncodeBase64(t *testing.T) {
  type args struct {
    data []byte
  }
  tests := []struct {
    name string
    args args
    want []byte
  }{
    {"test1", args{[]byte("abcdefg")}, []byte("YWJjZGVmZw==")},
    {"test2", args{[]byte("123456789")}, []byte("MTIzNDU2Nzg5")},
    {"test3", args{[]byte("faerferferf")}, []byte("ZmFlcmZlcmZlcmY=")},
    {"test4", args{[]byte("sde23rsdffvISibASOcvaKLS^*%^&$$%(&^R")}, []byte("c2RlMjNyc2RmZnZJU2liQVNPY3ZhS0xTXiolXiYkJCUoJl5S")},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := EncodeBase64(tt.args.data); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("EncodeBase64() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestDecodeBase64(t *testing.T) {
  type args struct {
    dataBase64 []byte
  }
  tests := []struct {
    name    string
    args    args
    want    []byte
    wantErr bool
  }{
    {"test1", args{[]byte("YWJjZGVmZw==")}, []byte("abcdefg"), false},
    {"test2", args{[]byte("MTIzNDU2Nzg5")}, []byte("123456789"), false},
    {"test3", args{[]byte("ZmFlcmZlcmZlcmY=")}, []byte("faerferferf"), false},
    {"test4", args{[]byte("c2RlMjNyc2RmZnZJU2liQVNPY3ZhS0xTXiolXiYkJCUoJl5S")}, []byte("sde23rsdffvISibASOcvaKLS^*%^&$$%(&^R"), false},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := DecodeBase64(tt.args.dataBase64)
      if (err != nil) != tt.wantErr {
        t.Errorf("DecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("DecodeBase64() = %v, want %v", got, tt.want)
      }
    })
  }
}
