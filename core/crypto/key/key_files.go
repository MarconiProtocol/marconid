package mcrypto_key

import (
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  mlog "github.com/MarconiProtocol/log"
  "io/ioutil"
  "os"
)

func savePrivateKey(filename string, key *rsa.PrivateKey) {
  keyFile, err := os.Create(filename)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to save key to file: ", filename, ", err = ", err)
  }
  defer keyFile.Close()

  var privateKey = &pem.Block{
    Type:  "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(key),
  }

  err = pem.Encode(keyFile, privateKey)
}

func savePublicKey(filename string, key *rsa.PublicKey) {
  ans1Bytes, err := x509.MarshalPKIXPublicKey(key)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to marshal public key", err)
  }

  var pemkey = &pem.Block{
    Type:  "PUBLIC KEY",
    Bytes: ans1Bytes,
  }

  keyFile, err := os.Create(filename)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to create file for key: ", filename, err)
  }
  defer keyFile.Close()

  err = pem.Encode(keyFile, pemkey)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to save key to file: ", filename, err)
  }
}

func readKeyFile(keyPath string) ([]byte, error) {
  return ioutil.ReadFile(keyPath)
}
