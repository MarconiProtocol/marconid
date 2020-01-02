package mutil

import (
  "crypto/rsa"
  "crypto/sha1"
  "crypto/x509"
  "encoding/base64"
  "encoding/hex"
  "encoding/pem"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "io/ioutil"
  "math/big"
  "os"
  "strconv"
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
  mask := (0xFFFFFFFF << (32 - uint(netmask))) & uint(0xFFFFFFFF)
  var dmask uint64
  dmask = 32
  localmask := ""
  for i := 1; i <= 4; i++ {
    tmp := mask >> (dmask - 8) & 0xFF

    if i == 1 {
      localmask += strconv.FormatUint(uint64(tmp), 10)
    } else {
      localmask += "." + strconv.FormatUint(uint64(tmp), 10)
    }
    dmask -= 8
  }
  return localmask
}

/*
  Return an integer that is deterministically calculated from two pubkeyhashes
  The resulting integer is used as the mpipe port between the two pubkeyhash owners as a form of psuedo port negotiation
*/
func GetMutualMPipePort(pubKeyHash string, peerPubKeyHash string) int {
  // choose a big enough prime number to use as the number of buckets
  bucketSize := big.NewInt(7919)
  // the base port is added to the result of the modulus to get the final port
  var basePort = 40000

  // concatenate the pubkeyhash strings based on a simple sorting order
  var concatPubKeyHash string
  if pubKeyHash >= peerPubKeyHash {
    concatPubKeyHash = pubKeyHash + peerPubKeyHash
  } else {
    concatPubKeyHash = peerPubKeyHash + pubKeyHash
  }

  // grab the bytes for the concatenated pubkey hashes and use the bytes to create a big int
  concatPubKeyHashBytes := []byte(concatPubKeyHash)
  num := &big.Int{}
  num.SetBytes(concatPubKeyHashBytes)
  // get the modulus of this big int number based on the prime bucket size
  res := &big.Int{}
  res.Mod(num, bucketSize)

  // The resulting port is the modulus added to the base port number
  return basePort + int(res.Int64())
}
