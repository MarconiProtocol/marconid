package msys_cmd_vars

import (
  "reflect"
  "testing"
)

func TestNewSafeWriter(t *testing.T) {
  tests := []struct {
    name string
    want *SafeWriter
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := NewSafeWriter(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("NewSafeWriter() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestSafeWriter_Write(t *testing.T) {
  type args struct {
    p []byte
  }
  tests := []struct {
    name    string
    sw      *SafeWriter
    args    args
    want    int
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := tt.sw.Write(tt.args.p)
      if (err != nil) != tt.wantErr {
        t.Errorf("SafeWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("SafeWriter.Write() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestSafeWriter_GetBytes(t *testing.T) {
  tests := []struct {
    name string
    sw   *SafeWriter
    want []byte
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.sw.GetBytes(); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("SafeWriter.GetBytes() = %v, want %v", got, tt.want)
      }
    })
  }
}
