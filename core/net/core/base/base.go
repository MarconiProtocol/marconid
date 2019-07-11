package mnet_core_base

import (
  "bytes"
  "encoding/binary"
  "errors"
  "hash"
  "time"

  "../../vars"

  "gitlab.neji.vm.tc/marconi/log"
)

var Log *mlog.Mlog

const (
  PIPE_MTU               = 1400 // mPipe's max-transmission-unit, largest data unit that can be transmitted over the pipes per transaction
  SOCKET_TIMEOUT_SECONDS = 3
)

/* Encode Frame
 * | HMAC-SHA256 (32 bytes) | Nanosecond Timestamp (8 bytes) |
 * |             Plaintext Frame (1-1432 bytes)              |
 */
func EncodeFrame(frame []byte, hmac_h hash.Hash) (enc_frame []byte, invalid error) {
  /* Encode Big Endian representation of current nanosecond unix time */
  time_unixnano := time.Now().UnixNano()
  time_bytes := make([]byte, 8)
  binary.BigEndian.PutUint64(time_bytes, uint64(time_unixnano))

  /* Prepend the timestamp to the frame */
  // timestamped frame + plaintext frame
  timestamped_frame := append(time_bytes, frame...)

  /* Compute the HMAC-SHA256 of the timestamped frame */
  hmac_h.Reset()
  hmac_h.Write(timestamped_frame)

  /* Prepend the HMAC-SHA256 */
  enc_frame = append(hmac_h.Sum(nil), timestamped_frame...)

  return enc_frame, nil
}

/*
   Decode Frame
*/
func DecodeFrame(enc_frame []byte, hmac_h hash.Hash) (frame []byte, invalid error) {
  /* Check that the encapsulated frame size is valid */
  if len(enc_frame) < (mnet_vars.TIMESTAMP_SIZE + mnet_vars.HMAC_SHA256_SIZE + 1) {
    return nil, errors.New("Invalid encapsulated frame size!")
  }

  /* Verify the timestamp */
  time_bytes := enc_frame[mnet_vars.HMAC_SHA256_SIZE : mnet_vars.HMAC_SHA256_SIZE+mnet_vars.TIMESTAMP_SIZE]
  time_unixnano := int64(binary.BigEndian.Uint64(time_bytes))
  curtime_unixnano := time.Now().UnixNano()
  if (curtime_unixnano - time_unixnano) > mnet_vars.TIMESTAMP_DIFF_THRESHOLD {
    return nil, errors.New("Timestamp outside of acceptable range!")
  }

  /* Verify the HMAC-SHA256 */
  //NOTE: if error out at here, possible, key does not match
  hmac_h.Reset()
  hmac_h.Write(enc_frame[mnet_vars.HMAC_SHA256_SIZE:])
  if bytes.Compare(hmac_h.Sum(nil), enc_frame[0:mnet_vars.HMAC_SHA256_SIZE]) != 0 {
    tmp := string(enc_frame[:])
    return nil, errors.New("Error verifying MAC! (Also check for key)" + tmp)
  }
  return enc_frame[mnet_vars.HMAC_SHA256_SIZE+mnet_vars.TIMESTAMP_SIZE:], nil
}
