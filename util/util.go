package mutil

import (
  "fmt"
  "os"
  "strconv"

  "crypto/rsa"
  "crypto/sha1"
  "crypto/x509"
  "encoding/base64"
  "encoding/hex"
  "encoding/pem"
  "gitlab.neji.vm.tc/marconi/log"
  "io/ioutil"
)


func DoesExist(filePath string) bool {
  if _, err := os.Stat(filePath); err == nil {
    return true
  }
  return false
}

func DecodeBase64(dataBase64 []byte) ([]byte, error) {
  /* Decode the base64 key */
  data := make([]byte, base64.StdEncoding.DecodedLen(len(dataBase64)))
  n, err := base64.StdEncoding.Decode(data, dataBase64)
  if err != nil {
    return nil, fmt.Errorf("Error decoding base64 key file: %s", err)
  }
  /* Truncate the key bytes to the right size */
  data = data[0:n]
  /* Check key size */
  if len(data) == 0 {
    return nil, fmt.Errorf("Error, invalid key in key file!")
  }
  return data, nil
}

func EncodeBase64(data []byte) []byte {
  /* Base64 encode the key */
  dataBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
  base64.StdEncoding.Encode(dataBase64, data)

  return dataBase64
}

func GetInfohashByPubKey(pub *rsa.PublicKey) (string, error) {
  pubKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
  if err != nil {
    mlog.GetLogger().Error("GetInfohashByPubKey failed due to error", err)
    return "", err
  }
  sha1Bytes := sha1.Sum(pubKeyBytes)
  return hex.EncodeToString(sha1Bytes[:]), nil
}

func LoadKey(privateKeyPath string) *rsa.PrivateKey {
  bytes, err := ioutil.ReadFile(privateKeyPath)
  if err != nil {
    mlog.GetLogger().Fatal("node_util: failed to open: ", privateKeyPath, err)
  }

  // decode PEM encoding to ANS.1 PKCS1 DER
  block, _ := pem.Decode(bytes)
  if block == nil {
    mlog.GetLogger().Fatal("node_util:LoadKey: No Block found in keyfile: ", privateKeyPath)
  }
  if block.Type != "RSA PRIVATE KEY" {
    mlog.GetLogger().Fatal("Unsupported key type")
  }
  // parse DER format to a native type
  key, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
  if err2 != nil {
    mlog.GetLogger().Fatalf("%v\n", err)
  }
  return key
}

/*
** helper function for converting a integer such as 24 net-mask to its 32-bit
** equivalent such as 255.255.255.0
 */
func Get32BitMaskFromCIDR(netmask int) string {
  mask := (0xFFFFFFFF << (32 - uint(netmask))) & 0xFFFFFFFF
  var dmask uint64
  dmask = 32
  localmask := ""
  for i := 1; i <= 4; i++ {
    tmp := mask >> (dmask - 8) & 0xFF

    if i == 1 {
      localmask += strconv.Itoa(tmp)
    } else {
      localmask += "." + strconv.Itoa(tmp)
    }
    dmask -= 8
  }
  return localmask
}
