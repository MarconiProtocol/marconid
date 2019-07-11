package mcrypto_key

import (
  "testing"
)

var km *KeyManager

func TestKeyManagerInstance(t *testing.T) {
  t.Log("TestKeyManagerInstance called")
  km = KeyManagerInstance()
  if km != KeyManagerInstance() {
    t.Error("can only have one KeyManager instance")
  }
}

func TestKeyManager_GetBasePrivateKeyAndPublicKey(t *testing.T) {
  t.Log("TestKeyManager_GetBasePrivateKey")
  km.EnsurePrivatePublicKeysGenerated()
  sk := km.GetBasePrivateKey()
  pk := km.GetBasePublicKey()

  if sk.N != pk.N || sk.E != pk.E {
    t.Error("public key and private key do not match")
  }
}
