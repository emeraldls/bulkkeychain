# bulkkeychain

`bulkkeychain` is an unofficial Go library for signing transactions for [Bulk](https://www.bulk.trade/).

> This library is not affiliated with or endorsed by Bulk. Use Bulk's [API documentation](https://docs.bulk.trade/api-reference/introduction) for transaction rules and request details.

## Install

```bash
go get github.com/emeraldls/bulkkeychain
```

## Sign a transaction

```go
keypair, err := bulkkeychain.NewKeyPair().WithBase58(base58PrivateKey)
if err != nil {
	return err
}
signer := bulkkeychain.NewSigner(keypair, bulkkeychain.Mainnet)

input := bulkkeychain.SignInput{
	Actions: []bulkkeychain.Action{
		{
			MarketOrder: &bulkkeychain.MarketOrderAction{
				Symbol: "BTC-USD",
				Buy:    true,
				Size:   0.1,
			},
		},
	},
	Nonce:   uint64(time.Now().UnixNano()),
	Account: accountPublicKey,
}

message, err := signer.Sign(input)
if err != nil {
	return err
}
```

`message` is a `SignMessage` ready to JSON-encode. The library signs the transaction but does not send it to Bulk.

## Serialize without signing

Use `SerializeMessage` when another service, such as Turnkey, holds the private key. Using the `input` above:

```go
signer := bulkkeychain.NewSigner(nil, bulkkeychain.Mainnet)
payload, err := signer.SerializeMessage(input)
if err != nil {
	return err
}

// Pass payload.Bytes() to your external signer.
```

The returned `bytes.Buffer` contains the actions, nonce, account, and network domain—the same bytes used by `Sign`. No private key is needed to serialize. Do not call `Sign` on a signer with a nil keypair.

## Key pairs

Generate a new key pair:

```go
keypair := bulkkeychain.NewKeyPair()
```

Load a Base58-encoded 64-byte Ed25519 private key:

```go
keypair := bulkkeychain.NewKeyPair().WithBase58(privateKey)
```

Get the raw keys with `keypair.PrivateKey()` and `keypair.PublicKey()`.

## Signer

Create one signer for the network you use:

```go
mainnetSigner := bulkkeychain.NewSigner(keypair, bulkkeychain.Mainnet)
testnetSigner := bulkkeychain.NewSigner(keypair, bulkkeychain.Testnet)
devnetSigner := bulkkeychain.NewSigner(keypair, bulkkeychain.Devnet)
```

The signer domain must match the Bulk network receiving the transaction.

## Actions

Each `Action` must contain exactly one action type. Add more `Action` values to sign a batch.

| Action | Go type |
| --- | --- |
| Market order | `MarketOrderAction` |
| Limit order | `LimitOrderAction` |
| Modify order | `ModifyOrderAction` |
| Cancel one order | `CancelSingleOrderAction` |
| Cancel all orders | `CancelAllOrdersAction` |
| Stop order | `StopOrderAction` |
| Take-profit order | `TakeProfitAction` |
| Transfer | `TransferAction` |
| Approve builder code | `BuilderCodeAction` |
| Revoke builder code | `RevokeBuilderCodeAction` |
| AgentWalletCreation | `AgentWalletCreationAction` |

Range/OCO, trigger basket, trailing stop, and on-fill actions are not supported.

## Links

- [Go reference](https://pkg.go.dev/github.com/emeraldls/bulkkeychain)
- [GitHub](https://github.com/emeraldls/bulkkeychain)
- [Bulk API](https://docs.bulk.trade/api-reference/introduction)
- [Bulk signing format](https://docs.bulk.trade/api-reference/signing)
