# bulkkeychain

[![Go Reference](https://pkg.go.dev/badge/github.com/emeraldls/bulkkeychain.svg)](https://pkg.go.dev/github.com/emeraldls/bulkkeychain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`bulkkeychain` is a Go library for serializing and signing order transactions for [Bulk](https://www.bulk.trade/).

> [!WARNING]
> This is an unofficial community library. It is not affiliated with or endorsed by Bulk, and it currently supports only the order actions listed below. Review and test signing behavior before using it with production funds.

The library was originally built to power [Clique](https://clique.trade), a social trading product where Telegram groups can share positions, copy friends, or countertrade them on Bulk.

## Installation

```bash
go get github.com/emeraldls/bulkkeychain
```

## Usage

Build an `OrderInput`, sign it with the private key belonging to the signer, then JSON-encode the returned `Order` for Bulk's `POST /order` endpoint.

```go
input := bulkkeychain.OrderInput{
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

order, err := bulkkeychain.SignOrder(input, privateKey, bulkkeychain.Mainnet)
if err != nil {
	return err
}

payload, err := json.Marshal(order)
if err != nil {
	return err
}

// Send payload to Bulk's POST /order endpoint.
```

`Account` is the Bulk account being traded. `SignOrder` derives `Signer` from `privateKey`, so an authorized agent wallet can sign for a different account. The library creates the signature but does not submit the request.

## Order support

| Order action | Tag | Support |
| --- | --- | --- |
| Market | `m` | Supported |
| Limit (GTC, IOC, ALO) | `l` | Supported |
| Modify | `mod` | Supported |
| Cancel one | `cx` | Supported |
| Cancel all | `cxa` | Supported |
| Stop | `st` | Supported |
| Take profit | `tp` | Supported |
| Range / OCO | `rng` | Not supported |
| Trigger basket | `trig` | Not supported |
| Trailing stop | `trl` | Not supported |
| On-fill | `of` | Not supported |

Market and limit orders also support optional builder-code payloads.

Additional order actions may be added based on community demand. Open an issue or pull request with the use case and the relevant Bulk protocol documentation.

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
