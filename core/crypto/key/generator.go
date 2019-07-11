package mcrypto_key

import (
  "../../../util"
  "crypto/rand"
  "crypto/rsa"
  "gitlab.neji.vm.tc/marconi/log"
)

func (km *KeyManager) generatePrivatePublicKeyPair() {
  random := rand.Reader
  privateKey, err := rsa.GenerateKey(random, 2048)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to generate private key for node")
  }
  savePrivateKey(PRIVATE_KEY_PATH, privateKey)
  savePublicKey(PUBLIC_KEY_PATH, &privateKey.PublicKey)

  pubKeyhash, _ := mutil.GetInfohashByPubKey(&privateKey.PublicKey)
  mlog.GetLogger().Info("\n===============")
  mlog.GetLogger().Info("Generated new private/public key pair, pubKeyhash: ", pubKeyhash)
  mlog.GetLogger().Info("===============\n")
}
