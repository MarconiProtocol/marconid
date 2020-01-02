package mnet_security

import (
  "bytes"
  "crypto/sha256"
  "encoding/binary"
  mlog "github.com/MarconiProtocol/log"
  "strings"
  "time"
)

const KEY_CHANGE_INTERVAL_SECONDS uint32 = 60

func GetBlockTimeIndex(intervalSeconds uint32) int64 {
  seconds := time.Now().Unix()
  // Adjust the time to a multiple of intervalSeconds.
  bucketed_seconds := (seconds / int64(intervalSeconds)) * int64(intervalSeconds)
  return bucketed_seconds
}

func KeepUpdateDataKey(dataKey *[]byte, keySignal *chan *bytes.Buffer, msgSignal chan string) {
  originalKey := make([]byte, len(*dataKey))
  copy(originalKey, *dataKey)
  for {
    select {
    case msg := <-msgSignal:
      if strings.Compare(msg, "quit") == 0 {
        return
      }
    case newKey := <-*keySignal:
      mlog.GetLogger().Debug("NewKey received and copied to original key")
      originalKey = make([]byte, len(*dataKey))
      copy(originalKey, newKey.Bytes())
    default:
      currentBlockTimeIndex := GetBlockTimeIndex(KEY_CHANGE_INTERVAL_SECONDS)
      UpdateTsDataKey(dataKey, originalKey, currentBlockTimeIndex)
      time.Sleep(3 * time.Second)
    }
  }
}

func UpdateTsDataKey(dataKey *[]byte, originalKey []byte, blockTime int64) []byte {
  hash_material := make([]byte, 40)
  copy(hash_material, originalKey)
  binary.LittleEndian.PutUint64(hash_material[32:], uint64(blockTime))
  sum := sha256.Sum256(hash_material)
  *dataKey = sum[:]
  ret := *dataKey
  return ret
}
