package bulkkeychain

import (
	"crypto/ed25519"

	"github.com/btcsuite/btcutil/base58"
)

func NewKeyPair() *KeyPair {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)

	k := &KeyPair{
		privateKey: privKey,
		publicKey:  pubKey,
	}
	return k
}

func (k *KeyPair) PrivateKey() ed25519.PrivateKey {
	return k.privateKey
}

func (k *KeyPair) PublicKey() ed25519.PublicKey {
	return k.publicKey
}

// Accepts Base58 encoded private key
func (k *KeyPair) WithBase58(secretKey string) *KeyPair {
	if k == nil {
		k = NewKeyPair()
	}

	decoded := base58.Decode(secretKey)
	if len(decoded) != ed25519.PrivateKeySize {
		panic("private key must be 64 bytes")
	}

	k.privateKey = ed25519.PrivateKey(decoded)
	k.publicKey = k.privateKey.Public().(ed25519.PublicKey)

	return k
}
