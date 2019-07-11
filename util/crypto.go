package mutil

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "errors"
  "io"
)

type SymmetricCrypter interface {
  Encrypt(plaintext []byte) (ciphertext []byte, err error)
  Decrypt(ciphertext []byte) (plaintext []byte, err error)
}

type aesCrypter struct {
  key []byte
}

func NewAesCrypter(key []byte) *aesCrypter {
  crypter := new(aesCrypter)
  crypter.key = key
  return crypter
}

func (ac *aesCrypter) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
  var block cipher.Block

  if block, err = aes.NewCipher(ac.key); err != nil {
    return nil, err
  }

  ciphertext = make([]byte, aes.BlockSize+len(plaintext))

  initialization_vector := ciphertext[:aes.BlockSize]
  if _, err = io.ReadFull(rand.Reader, initialization_vector); err != nil {
    return nil, err
  }

  cfb := cipher.NewCFBEncrypter(block, initialization_vector)
  cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
  return ciphertext, nil
}

func (ac *aesCrypter) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
  var block cipher.Block

  if block, err = aes.NewCipher(ac.key); err != nil {
    return nil, err
  }

  if len(ciphertext) < aes.BlockSize {
    err = errors.New("ciphertext too short")
    return nil, err
  }

  initialization_vector := ciphertext[:aes.BlockSize]
  ciphertext = ciphertext[aes.BlockSize:]

  cfb := cipher.NewCFBDecrypter(block, initialization_vector)
  cfb.XORKeyStream(ciphertext, ciphertext)
  return ciphertext, nil
}
