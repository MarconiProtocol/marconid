package mutil

import (
  "bytes"
  "testing"
)

func getText() []byte {
  return []byte("some text")
}

func TestAes128(t *testing.T) {
  key_16_bytes := []byte("deadbeefdeadbeef")
  var crypter SymmetricCrypter = NewAesCrypter(key_16_bytes)
  ciphertext, err := crypter.Encrypt(getText())
  if err != nil {
    t.Error("Unexpected encryption error:", err)
  }
  plaintext, err := crypter.Decrypt(ciphertext)
  if err != nil {
    t.Error("Unexpected decryption error:", err)
  }
  if !bytes.Equal(plaintext, getText()) {
    t.Error("Incorrect decrypted text:", string(plaintext))
  }
}

func TestAes256(t *testing.T) {
  key_32_bytes := []byte("deadbeefdeadbeefdeadbeefdeadbeef")
  var crypter SymmetricCrypter = NewAesCrypter(key_32_bytes)
  ciphertext, err := crypter.Encrypt(getText())
  if err != nil {
    t.Error("Unexpected encryption error:", err)
  }
  plaintext, err := crypter.Decrypt(ciphertext)
  if err != nil {
    t.Error("Unexpected decryption error:", err)
  }
  if !bytes.Equal(plaintext, getText()) {
    t.Error("Incorrect decrypted text:", string(plaintext))
  }
}

func TestBadKey(t *testing.T) {
  key_17_bytes := []byte("deadbeefdeadbeefz")
  var crypter SymmetricCrypter = NewAesCrypter(key_17_bytes)
  _, err := crypter.Encrypt(getText())
  if err == nil {
    t.Error("Expected a key size error but got none.")
  }
}

func TestAesMismatchedKeys(t *testing.T) {
  var crypter SymmetricCrypter
  encryption_key := []byte("deadbeefdeadbeef")
  decryption_key := []byte("beefdeadbeefdead")

  crypter = NewAesCrypter(encryption_key)
  ciphertext, err := crypter.Encrypt(getText())
  if err != nil {
    t.Error("Unexpected encryption error:", err)
  }
  crypter = NewAesCrypter(decryption_key)
  plaintext, err := crypter.Decrypt(ciphertext)
  if err != nil {
    t.Error("Unexpected decryption error:", err)
  }
  if bytes.Equal(plaintext, getText()) {
    t.Error("Expected incorrect decrypted text.")
  }
}
