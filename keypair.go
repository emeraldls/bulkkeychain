package bulkkeychain

import (
	"crypto/ed25519"
	"errors"

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

// For existing wallet
func (k *KeyPair) WithBase58(secretKey string) (*KeyPair, error) {
	if k == nil {
		k = NewKeyPair()
	}

	decoded := base58.Decode(secretKey)
	if len(decoded) != ed25519.PrivateKeySize {
		return nil, errors.New("private key must be 64 bytes")
	}

	k.privateKey = ed25519.PrivateKey(decoded)
	k.publicKey = k.privateKey.Public().(ed25519.PublicKey)

	return k, nil
}
