package mutil

import (
  "../core/sys/cmd"
  "fmt"
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
