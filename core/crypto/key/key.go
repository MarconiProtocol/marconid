package mcrypto_key

import (
  "../../net/vars"
  "bytes"
  "crypto/rand"
  "encoding/base64"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "io/ioutil"
  "os"
)

func GetLinkLayerEncryptionKey(privateKeyFilePath string, generateIfNotExist bool) []byte {
  //var key []byte
  key, err := Keyfile_read(privateKeyFilePath)
  /* If the error is file does not exist */
  if err != nil && os.IsNotExist(err) {
    if generateIfNotExist {
      /* Auto-generate the key file */
      key, err = Keyfile_generate(privateKeyFilePath)
      if err != nil {
        mlog.GetLogger().Fatalf("Error generating key file: %s", err)
      }
    }
  } else if err != nil {
    mlog.GetLogger().Fatalf("Error reading key file: %s", err)
  }
  return key
}

func GenerateKey() []byte {
  key := make([]byte, mnet_vars.DATA_KEY_SIZE)
  //rand.Seed(int64(seed))
  _, err := rand.Read(key)
  if err != nil {
    fmt.Println("Error random key:", key)
  }
  fmt.Println("key in string: ", string(key[:]), len(string(key[:])))
  fmt.Println("key in byte: ", key, "len: ", len(key))

  //  keyStr := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
  //  base64.StdEncoding.Encode(keyStr, key)

  keyStr := base64.StdEncoding.EncodeToString(key)
  fmt.Println("keyStr: ", keyStr, "[]byte(keyStr):", []byte(keyStr))

  saveKey("../build/test_key", []byte(keyStr))
  //return string(keyStr[:])
  return key
}

func saveKey(filename string, content []byte) {
  err := ioutil.WriteFile(filename, content, 0644)
  if err != nil {
    fmt.Println("Error saving key file:", filename, content, "err =", err)
  }
}

/**********************************************************************/
/*** Key file reading and generation ***/
/**********************************************************************/

/* The key file simply contains a base64 encoded random key.
 * The default random key size is HMAC_SHA256_SIZE. */
func Keyfile_read(path string) (key []byte, e error) {
  var key_base64 []byte

  /* Attempt to open the key file for reading */
  keyfile, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer keyfile.Close()

  /* Get the key file size */
  fi, err := keyfile.Stat()
  if err != nil {
    return nil, fmt.Errorf("Error getting key file size: %s", err)
  }

  /* Read the base64 key */
  key_base64 = make([]byte, fi.Size())
  n, err := keyfile.Read(key_base64)
  if err != nil {
    return nil, fmt.Errorf("Error reading key file: %s", err)
  }
  /* Trim whitespace */
  key_base64 = bytes.TrimSpace(key_base64)

  /* Decode the base64 key */
  key = make([]byte, base64.StdEncoding.DecodedLen(len(key_base64)))
  n, err = base64.StdEncoding.Decode(key, key_base64)
  if err != nil {
    return nil, fmt.Errorf("Error decoding base64 key file: %s", err)
  }
  /* Truncate the key bytes to the right size */
  key = key[0:n]

  /* Check key size */
  if len(key) == 0 {
    return nil, fmt.Errorf("Error, invalid key in key file!")
  }
  return key, nil
}

func Keyfile_generate(path string) (key []byte, e error) {
  /* Generate a random key */
  key = make([]byte, mnet_vars.HMAC_SHA256_SIZE)
  n, err := rand.Read(key)
  if n != len(key) {
    return nil, fmt.Errorf("Error generating random key of size %d: %s", len(key), err)
  }

  /* Base64 encode the key */
  key_base64 := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
  base64.StdEncoding.Encode(key_base64, key)

  /* Open the key file for writing */
  keyfile, err := os.Create(path)
  if err != nil {
    return nil, fmt.Errorf("Error opening key file for writing: %s", err)
  }
  defer keyfile.Close()

  /* Write the base64 encoded key */
  _, err = keyfile.Write(key_base64)
  if err != nil {
    return nil, fmt.Errorf("Error writing base64 encoded key to keyfile: %s", err)
  }

  return key, nil
}
