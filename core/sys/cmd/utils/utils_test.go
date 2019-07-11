package msys_cmd_utils

import (
  "os/exec"
  "testing"
)

func TestParseLinuxVersion(t *testing.T) {
  tests := []struct {
    name    string
    want    int
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := ParseLinuxVersion()
      if (err != nil) != tt.wantErr {
        t.Errorf("ParseLinuxVersion() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("ParseLinuxVersion() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestRunCmd(t *testing.T) {
  type args struct {
    cmd *exec.Cmd
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := RunCmd(tt.args.cmd)
      if (err != nil) != tt.wantErr {
        t.Errorf("RunCmd() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("RunCmd() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestRunPipedCmds(t *testing.T) {
  type args struct {
    cmds []*exec.Cmd
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := RunPipedCmds(tt.args.cmds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("RunPipedCmds() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("RunPipedCmds() = %v, want %v", got, tt.want)
      }
    })
  }
}
