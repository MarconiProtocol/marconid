package mcrypto_key

import (
  "bytes"
  "os"
  "testing"
)

var keyPath = "../build/test_key"

func TestGenerateAndReadKey(t *testing.T) {
  t.Log("TestGenerateAndReadKey called")
  key := GenerateKey()
  if key == nil {
    t.Error("GenerateKey failed")
  }

  key1, err := Keyfile_read(keyPath)
  if err != nil || bytes.Compare(key, key1) != 0 {
    t.Error("Keyfile_read failed")
  }
}

func TestGetLinkLayerEncryptionKey(t *testing.T) {
  t.Log("TestGetLinkLayerEncryptionKey called")
  // remove old key file
  _, err := os.Stat(keyPath)
  if err == nil {
    os.Remove(keyPath)
  }

  key := GetLinkLayerEncryptionKey(keyPath, false)
  if key != nil {
    t.Error("should not generate keyFile if does not exist")
  }

  key = GetLinkLayerEncryptionKey(keyPath, true)
  _, err = os.Stat(keyPath)
  key1, err := Keyfile_read(keyPath)
  if bytes.Compare(key, key1) != 0 {
    t.Error("GetLinkLayerEncryptionKey failed")
  }
}
