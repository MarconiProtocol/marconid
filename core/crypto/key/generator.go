package mcrypto_key

import (
  "../../../util"
  "../../config"
  "crypto/rand"
  "crypto/rsa"
  mlog "github.com/MarconiProtocol/log"
)

func (km *KeyManager) generatePrivatePublicKeyPair() {
  random := rand.Reader
  privateKey, err := rsa.GenerateKey(random, 2048)
  if err != nil {
    mlog.GetLogger().Fatal("Failed to generate private key for node: ", err)
  }
  savePrivateKey(mconfig.GetUserConfig().Global.Base_Dir+PRIVATE_KEY_PATH, privateKey)
  savePublicKey(mconfig.GetUserConfig().Global.Base_Dir+PUBLIC_KEY_PATH, &privateKey.PublicKey)

  pubKeyhash, _ := mutil.GetInfohashByPubKey(&privateKey.PublicKey)
  mlog.GetLogger().Info("Generated new private/public key pair, pubKeyhash=", pubKeyhash)
}
