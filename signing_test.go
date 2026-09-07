package bulkkeychain

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

func TestSign(t *testing.T) {
	pubKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal()
	}

	account := base58.Encode([]byte("gurtcryptogurtcryptogurtcryptogu"))
	nonce := 123456789

	signInput := SignInput{
		Actions: []Action{
			{
				MarketOrder: &MarketOrderAction{
					Symbol:          "ETH-USD",
					Buy:             false,
					Size:            0.1,
					ReduceOnly:      false,
					IsolatedAccount: true,
					BuilderCode:     nil,
				},
			},
			{
				LimitOrder: &LimitOrderAction{

					Symbol:          "ETH-USD",
					Buy:             false,
					Size:            0.1,
					ReduceOnly:      false,
					IsolatedAccount: true,
					BuilderCode:     nil,
				},
			},
		},
		Nonce:   uint64(nonce),
		Account: account,
	}

	keypair := NewKeyPair().WithBase58(base58.Encode(privateKey))

	signer := NewSigner(keypair, Mainnet)

	message, err := signer.Sign(signInput)
	if err != nil {
		t.Fatal(err)
	}

	expectedSigner := base58.Encode(pubKey)

	if message.Signer != expectedSigner {
		t.Fatalf("exptected signer: %s, got: %s", expectedSigner, message.Signer)
	}

	if message.Account != account {
		t.Fatalf("expected account: %s, got %s", account, message.Account)
	}

	if message.Nonce != uint64(nonce) {
		t.Fatalf("expected none: %d, got %d", nonce, message.Nonce)
	}

	signature := base58.Decode(message.Signature)
	if len(signature) != ed25519.SignatureSize {
		t.Fatalf("expected signature size: %d, got: %d", ed25519.SignatureSize, len(signature))
	}

	var expectedMessage bytes.Buffer
	actions, err := serializeActions(signInput.Actions)
	if err != nil {
		t.Fatal(err)
	}
	expectedMessage.Write(actions)
	if err := binary.Write(&expectedMessage, binary.LittleEndian, signInput.Nonce); err != nil {
		t.Fatal(err)
	}
	expectedMessage.Write(base58.Decode(account))
	expectedMessage.WriteByte(byte(Mainnet))

	if !ed25519.Verify(pubKey, expectedMessage.Bytes(), signature) {
		t.Fatal("unable to verify signature")
	}
}

func TestSignRejectsInvalidPrivateKey(t *testing.T) {
	signer := NewSigner(&KeyPair{privateKey: ed25519.PrivateKey("too short")}, Mainnet)
	_, err := signer.Sign(SignInput{})
	if err == nil {
		t.Fatal("expected invalid private key error")
	}
}

func TestSerializeActionsRequiresExactlyOneOrderType(t *testing.T) {
	tests := []struct {
		name   string
		action Action
	}{
		{
			name:   "empty",
			action: Action{},
		},
		{
			name: "multiple",
			action: Action{
				MarketOrder: &MarketOrderAction{},
				LimitOrder:  &LimitOrderAction{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := serializeActions([]Action{test.action})
			if err == nil {
				t.Fatal("expected invalid action error")
			}
		})
	}
}
