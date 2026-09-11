package bulkkeychain

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

var benchmarkKeyPair *KeyPair
var benchmarkMessage *SignMessage

func BenchmarkNewKeyPair(b *testing.B) {
	for b.Loop() {
		benchmarkKeyPair = NewKeyPair()
	}
}

func BenchmarkSign(b *testing.B) {
	keypair := NewKeyPair()
	signer := NewSigner(keypair, Devnet)
	account := base58.Encode(keypair.PublicKey())

	for _, actionCount := range []int{1, 10, 100, 1_000_000} {
		b.Run(fmt.Sprintf("%d_actions", actionCount), func(b *testing.B) {
			action := &MarketOrderAction{
				Symbol: "BTC-USD",
				Buy:    true,
				Size:   0.1,
			}
			actions := make([]Action, actionCount)
			for i := range actions {
				actions[i].MarketOrder = action
			}
			input := SignInput{
				Actions: actions,
				Nonce:   1234567890,
				Account: account,
			}

			b.ReportAllocs()
			b.ReportMetric(float64(actionCount), "actions/op")
			b.ResetTimer()
			for b.Loop() {
				message, err := signer.Sign(input)
				if err != nil {
					b.Fatal(err)
				}
				benchmarkMessage = message
			}
		})
	}
}
