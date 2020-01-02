package msys_cmd

import (
  "./vars"
  "bytes"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "io"
  "os"
  "os/exec"
)

func ExecuteCommand(cmd string, cmdArgs []string) (result string) {

  fmt.Println(exec.Command("echo", "$PATH"))

  fmt.Println("cmd #: ", cmd, " cmdArgs: ", cmdArgs)
  out, err := exec.Command(cmd, cmdArgs...).Output()
  if err != nil {
    fmt.Println("ExecuteCommand: error: ", err)
    //log.Fatal(err)
  }
  result = string(out)
  fmt.Println("out: ", result)
  return
}

func ExecuteCommandPipe(cmd1 string, cmd1Args []string, cmd2 string, cmd2Args []string) (result bytes.Buffer) {
  c1 := exec.Command(cmd1, cmd1Args...)
  c2 := exec.Command(cmd2, cmd2Args...)

  r, w := io.Pipe()
  c1.Stdout = w
  c2.Stdin = r

  c2.Stdout = &result

  c1.Start()
  c2.Start()
  c1.Wait()
  w.Close()
  c2.Wait()
  io.Copy(os.Stdout, &result)
  return
}

func ExecuteSequencialIdenticalCommand(cmd string, cmdArgsList map[int][]string) (result map[int]string) {
  result = make(map[int]string)
  for index, cmdArgs := range cmdArgsList {
    result[index] = ExecuteCommand(cmd, cmdArgs)
  }
  return
}

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
  //TODO: check for cmd/cmd args
  // Collect the output from the command(s)
  output := msys_cmd_vars.NewSafeWriter()
  stderr := msys_cmd_vars.NewSafeWriter()

  last := len(cmds) - 1
  for i, cmd := range cmds[:last] {
    var err error
    // Connect each command's stdin to the previous command's stdout
    if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
      return nil, nil, err
    }
    // Connect each command's stderr to a buffer
    cmd.Stderr = stderr
  }

  // Connect the output and error for the last command
  cmds[last].Stdout, cmds[last].Stderr = output, stderr

  // Start each command
  for _, cmd := range cmds {
    mlog.GetLogger().Debug("starting the cmd", cmd.Args)
    if err := cmd.Start(); err != nil {
      return output.GetBytes(), stderr.GetBytes(), err
    }
  }

  // Wait for each command to complete
  for _, cmd := range cmds {
    mlog.GetLogger().Debug("waiting on the cmd", cmd.Args)
    if err := cmd.Wait(); err != nil {
      return output.GetBytes(), stderr.GetBytes(), err
    }
  }

  // Return the pipeline output and the collected standard error
  return output.GetBytes(), stderr.GetBytes(), nil
}
