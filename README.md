# bulkkeychain

[![Go Reference](https://pkg.go.dev/badge/github.com/emeraldls/bulkkeychain.svg)](https://pkg.go.dev/github.com/emeraldls/bulkkeychain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`bulkkeychain` is an unofficial Go library for serializing and signing transactions for [Bulk](https://www.bulk.trade/).

> [!WARNING]
> This is an unofficial community library. It is not affiliated with or endorsed by Bulk, and it currently supports only the actions listed below. Review and test signing behavior before using it with production funds.

I started building this library to power [Clique](https://clique.trade) because I needed Bulk transaction signing in Go. Clique is a social trading product where friends in Telegram groups can share positions, copy friends, or countertrade them on Bulk.
## Installation

```bash
go get github.com/emeraldls/bulkkeychain
```

## Usage

Create a key pair from your Base58 private key, create a signer, then sign a `SignInput`. JSON-encode the returned `SignMessage` for Bulk's `POST /order` endpoint.

```go
keypair := bulkkeychain.NewKeyPair().WithBase58(base58PrivateKey)
signer := bulkkeychain.NewSigner(keypair, bulkkeychain.Mainnet)

input := bulkkeychain.SignInput{
	Actions: []bulkkeychain.Action{
		{
			MarketOrder: &bulkkeychain.MarketOrderAction{
				Symbol:          "BTC-USD",
				Buy:             true,
				Size:            0.1,
				ReduceOnly:      false,
				IsolatedAccount: false,
			},
		},
	},
	Nonce:   uint64(time.Now().UnixNano()),
	Account: accountPublicKey,
}

signed, err := signer.Sign(input)
if err != nil {
	return err
}

payload, err := json.Marshal(signed)
if err != nil {
	return err
}

// Send payload to Bulk's POST /order endpoint.
```

`Account` is the Bulk account authorizing the transaction. The signer public key is derived from the key pair, so an authorized agent wallet can sign for a different account. The library creates the signature but does not submit the request.

## Action support

| Action | Tag | Support |
| --- | --- | --- |
| Market | `m` | Supported |
| Limit (GTC, IOC, ALO) | `l` | Supported |
| Modify | `mod` | Supported |
| Cancel one | `cx` | Supported |
| Cancel all | `cxa` | Supported |
| Stop | `st` | Supported |
| Take profit | `tp` | Supported |
| Transfer | `transfer` | Supported |
| Approve builder code | `abc` | Supported |
| Revoke builder code | `rbc` | Supported |

Market and limit orders also support optional builder-code payloads.

Additional actions may be added based on community demand. Open an issue or pull request with the use case and the relevant Bulk protocol documentation.

## Signing format

Signatures follow Bulk's canonical binary format:

```text
wincode(actions) || nonce_le_u64 || account_pubkey || signature_domain
```

The supported domains are `Mainnet`, `Testnet`, and `Devnet`. See Bulk's [transaction signing documentation](https://docs.bulk.trade/api-reference/signing) and [order API reference](https://docs.bulk.trade/api-reference/placeOrder) for the protocol specification.

## Development

```bash
go test ./...
```

## License

[MIT](LICENSE)
