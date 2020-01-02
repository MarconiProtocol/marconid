package mcrypto_dh

import (
  "bytes"
  "github.com/MarconiProtocol/kyber"
  "testing"
)

func TestDHKeyInfo_EncodePublicKeyToString(t *testing.T) {
  type fields struct {
    PrivateKey        kyber.Scalar
    PublicKey         kyber.Point
    PeerPublicKey     kyber.Point
    SymmetricKey      kyber.Point
    SymmetricKeyBytes *bytes.Buffer
    KeySignal         *chan *bytes.Buffer
  }
  tests := []struct {
    name   string
    fields fields
    want   string
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      dhKeyInfo := &DHKeyInfo{
        PrivateKey:        tt.fields.PrivateKey,
        PublicKey:         tt.fields.PublicKey,
        PeerPublicKey:     tt.fields.PeerPublicKey,
        SymmetricKey:      tt.fields.SymmetricKey,
        SymmetricKeyBytes: tt.fields.SymmetricKeyBytes,
        KeySignal:         tt.fields.KeySignal,
      }
      if got := dhKeyInfo.EncodePublicKeyToString(); got != tt.want {
        t.Errorf("DHKeyInfo.EncodePublicKeyToString() = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestDHKeyInfo_DecodeStringToPeerPublicKey(t *testing.T) {
  type fields struct {
    PrivateKey        kyber.Scalar
    PublicKey         kyber.Point
    PeerPublicKey     kyber.Point
    SymmetricKey      kyber.Point
    SymmetricKeyBytes *bytes.Buffer
    KeySignal         *chan *bytes.Buffer
  }
  type args struct {
    encodedPublicKey string
  }
  tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      dhKeyInfo := &DHKeyInfo{
        PrivateKey:        tt.fields.PrivateKey,
        PublicKey:         tt.fields.PublicKey,
        PeerPublicKey:     tt.fields.PeerPublicKey,
        SymmetricKey:      tt.fields.SymmetricKey,
        SymmetricKeyBytes: tt.fields.SymmetricKeyBytes,
        KeySignal:         tt.fields.KeySignal,
      }
      if err := dhKeyInfo.DecodeStringToPeerPublicKey(tt.args.encodedPublicKey); (err != nil) != tt.wantErr {
        t.Errorf("DHKeyInfo.DecodeStringToPeerPublicKey() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestDHKeyInfo_GenDHSymmetricKey(t *testing.T) {
  type fields struct {
    PrivateKey        kyber.Scalar
    PublicKey         kyber.Point
    PeerPublicKey     kyber.Point
    SymmetricKey      kyber.Point
    SymmetricKeyBytes *bytes.Buffer
    KeySignal         *chan *bytes.Buffer
  }
  tests := []struct {
    name   string
    fields fields
  }{
  // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      dhKeyInfo := &DHKeyInfo{
        PrivateKey:        tt.fields.PrivateKey,
        PublicKey:         tt.fields.PublicKey,
        PeerPublicKey:     tt.fields.PeerPublicKey,
        SymmetricKey:      tt.fields.SymmetricKey,
        SymmetricKeyBytes: tt.fields.SymmetricKeyBytes,
        KeySignal:         tt.fields.KeySignal,
      }
      dhKeyInfo.GenDHSymmetricKey()
    })
  }
}
