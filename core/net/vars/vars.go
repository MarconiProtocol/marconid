package mnet_vars

import (
  "crypto/sha256"
)

const (
  MPIPE_LINK_DEVICE_NAME_PREFIX    = "mp"  //l2 mpipe
  MPIPE_LINK_L3_DEVICE_NAME_PREFIX = "mpl" //l3 mpipe light
  MBINDER_LINK_DEVICE_NAME_PREFIX  = "mb"  //mbridge
)

const EMPTY_CONTRACT_ADDRESS = "0x0000000000000000000000000000000000000000"

const (
  DATA_KEY_SIZE = 32
  HMAC_SHA256_SIZE = sha256.Size
  TIMESTAMP_SIZE = 8
  TIMESTAMP_DIFF_THRESHOLD int64 = 10.0 * 1e9   /* Acceptable timestamp difference threshold in nS (10.0 seconds) */
  UDP_MTU = 1472
  TAP_MTU = UDP_MTU - 14 - HMAC_SHA256_SIZE - TIMESTAMP_SIZE - DATA_KEY_SIZE
  TUN_MTU = UDP_MTU - 20 - 20 - HMAC_SHA256_SIZE - TIMESTAMP_SIZE - DATA_KEY_SIZE
)
