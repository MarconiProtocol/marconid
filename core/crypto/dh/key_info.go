package mcrypto_dh

import (
  "bytes"
  "encoding/base64"
  "gitlab.neji.vm.tc/marconi/kyber"
  "gitlab.neji.vm.tc/marconi/log"
)

type DHKeyInfo struct {
  PrivateKey        kyber.Scalar
  PublicKey         kyber.Point
  PeerPublicKey     kyber.Point
  SymmetricKey      kyber.Point
  SymmetricKeyBytes *bytes.Buffer

  KeySignal *chan *bytes.Buffer
}

func (dhKeyInfo *DHKeyInfo) EncodePublicKeyToString() string {
  buffer := new(bytes.Buffer)
  dhKeyInfo.PublicKey.MarshalTo(buffer)
  encodedPublicKey := base64.StdEncoding.EncodeToString(buffer.Bytes())
  return encodedPublicKey
}

func (dhKeyInfo *DHKeyInfo) DecodeStringToPeerPublicKey(encodedPublicKey string) error {
  decodedBytes, b64decodeErr := base64.StdEncoding.DecodeString(encodedPublicKey)
  if b64decodeErr != nil {
    mlog.GetLogger().Errorf("Error base64 decoding the string payload : %v", encodedPublicKey)
    return b64decodeErr
  }

  buffer := bytes.NewReader(decodedBytes)
  publicKey := DHExchangeManagerInstance().suite.Point()
  _, unmarshallErr := publicKey.UnmarshalFrom(buffer)
  if unmarshallErr != nil {
    mlog.GetLogger().Error("Error unmarshalling bytes from buffer into Point obj")
    return unmarshallErr
  }

  dhKeyInfo.PeerPublicKey = publicKey

  return nil
}

func (dhKeyInfo *DHKeyInfo) GenDHSymmetricKey() {
  dhKeyInfo.SymmetricKey = DHExchangeManagerInstance().suite.Point().Mul(dhKeyInfo.PrivateKey, dhKeyInfo.PeerPublicKey)
  dhKeyInfo.SymmetricKeyBytes.Reset()
  dhKeyInfo.SymmetricKey.MarshalTo(dhKeyInfo.SymmetricKeyBytes)
  // NOTE: we are using the goroutine system as a pseudo buffer to avoid blocking...
  // This guarentees delivery without blocking the main goroutine
  go func() {
    *dhKeyInfo.KeySignal <- dhKeyInfo.SymmetricKeyBytes
    mlog.GetLogger().Debug("New Symmetric Key, pushed to KeySignal")
  }()
}
