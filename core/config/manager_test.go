package mconfig

import "testing"

func TestInitializeConfigs(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      InitializeConfigs("/opt/marconi")
    })
  }
}

func TestInitializeAppConfig(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      InitializeAppConfig("/opt/marconi")
    })
  }
}

func Test_readAndLoadAppConfig(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      readAndLoadAppConfig()
    })
  }
}

func TestInitializeUserConfig(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      InitializeUserConfig("/opt/marconi")
    })
  }
}

func Test_readAndLoadUserConfig(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      readAndLoadUserConfig()
    })
  }
}

func Test_createDefaultUserConfigFile(t *testing.T) {
  tests := []struct {
    name    string
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := createDefaultUserConfigFile(); (err != nil) != tt.wantErr {
        t.Errorf("createDefaultUserConfigFile() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func Test_setAppConfigDefaults(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      setAppConfigDefaults()
    })
  }
}

func Test_setUserConfigDefaults(t *testing.T) {
  tests := []struct {
    name string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      setUserConfigDefaults()
    })
  }
}
