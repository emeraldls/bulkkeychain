package bulkkeychain

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

func TestKeyPair(t *testing.T) {
	original := NewKeyPair()
	loaded, err := NewKeyPair().WithBase58(base58.Encode(original.PrivateKey()))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(loaded.PrivateKey(), original.PrivateKey()) {
		t.Fatal("private keys do not match")
	}
	if !bytes.Equal(loaded.PublicKey(), original.PublicKey()) {
		t.Fatal("public keys do not match")
	}
}
