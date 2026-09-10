# bulkkeychain

[![Go Reference](https://pkg.go.dev/badge/github.com/emeraldls/bulkkeychain.svg)](https://pkg.go.dev/github.com/emeraldls/bulkkeychain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`bulkkeychain` is an unofficial Go library for serializing and signing transactions for [Bulk](https://www.bulk.trade/).

> This is an unofficial community library. It is not affiliated with or endorsed by Bulk, and it currently supports only the actions listed below. 

I started building this library to power [Clique](https://clique.trade) because I needed Bulk transaction signing in Go. Clique is a social trading product where friends in Telegram groups can share positions, copy friends, or countertrade them on Bulk.

## Installation

```bash
go get github.com/emeraldls/bulkkeychain
```

## Usage

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

I only supported the actions I needed for Clique, but I can add more based on community requests & I welcome contributions for additional actions. The following table lists the supported actions and their Bulk API tags.

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
| Agent Wallet Creation | `agentWalletCreation` | Supported |

## Signing format

Signatures follow Bulk's canonical binary format:

```text
wincode(actions) || nonce_le_u64 || account_pubkey || signature_domain
```

The supported domains are `Mainnet`, `Testnet`, and `Devnet`. See Bulk's [transaction signing documentation](https://docs.bulk.trade/api-reference/signing) and [order API reference](https://docs.bulk.trade/api-reference/placeOrder) for the protocol specification.

## Benchmarks

Bulk's Rust implementation vs mine, benchmarked on macOS 26.6 with an Apple M1 Pro, Go 1.26.4, and Rust 1.95.0.

| Operation | Mine (Go) | Bulk (Rust) | Faster |
| --- | ---: | ---: | --- |
| Generate key pair | 16.1 µs | 13.1 µs | Bulk 1.23× |
| Sign 1 action | 27.4 µs | 22.9 µs | Bulk 1.20× |
| Sign 10 actions | 29.2 µs | 28.3 µs | About equal |
| Sign 100 actions | 42.4 µs | 77.5 µs | Mine 1.83× |
| Sign 1,000,000 actions | 134 ms | 682 ms | Mine 5.08× |

The benchmarks use each library's public signing API with input creation excluded. Order ID generation was disabled in Rust. The million-action case is only a stress test.

```bash
go test -run='^$' -bench=. -benchmem
```

## Development

```bash
go test ./...
```

## License

[MIT](LICENSE)
