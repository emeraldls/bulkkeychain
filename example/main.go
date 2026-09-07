package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/emeraldls/bulkkeychain"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env: ", err)
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY is missing from .env")
	}

	keypair := bulkkeychain.NewKeyPair().WithBase58(privateKey)
	signer := bulkkeychain.NewSigner(keypair, bulkkeychain.Testnet)

	message, err := signer.Sign(bulkkeychain.SignInput{
		Actions: []bulkkeychain.Action{
			{
				MarketOrder: &bulkkeychain.MarketOrderAction{
					Symbol: "BTC-USD",
					Buy:    true,
					Size:   0.001,
				},
			},
		},
		Nonce:   uint64(time.Now().UnixNano()),
		Account: base58.Encode(keypair.PublicKey()),
	})
	if err != nil {
		log.Fatal("sign transaction: ", err)
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Fatal("encode request: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "https://exchange-api.bulk.trade/api/v1/order", bytes.NewReader(body))
	if err != nil {
		log.Fatal("create request: ", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 15 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("send request: ", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("read response: ", err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, responseBody, "", "  "); err != nil {
		log.Fatal("format response: ", err)
	}

	fmt.Println(response.Status)
	fmt.Println(pretty.String())
}
