package mcore

const (
  TYPE_OS_LINUX   = 0
  TYPE_OS_DARWIN  = 1
  TYPE_OS_WINDOWS = 2
  TYPE_OS_ANDROID = 3
  TYPE_OS_IOS     = 4
  TYPE_OS_UNKNOWN = 5
)

var OSStringToInt = map[string]int{
  "linux":   TYPE_OS_LINUX,
  "darwin":  TYPE_OS_DARWIN,
  "windows": TYPE_OS_WINDOWS,
  "android": TYPE_OS_ANDROID,
  "ios":     TYPE_OS_IOS,
}
