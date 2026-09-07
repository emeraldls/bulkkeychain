package bulkkeychain

import (
	"crypto/ed25519"
)

type Signing interface {
	Discriminant() uint32
}

type Action struct {
	LimitOrder         *LimitOrderAction        `json:"l,omitempty"`
	MarketOrder        *MarketOrderAction       `json:"m,omitempty"`
	ModifyOrder        *ModifyOrderAction       `json:"mod,omitempty"`
	CancelOrder        *CancelSingleOrderAction `json:"cx,omitempty"`
	CancelAllOrders    *CancelAllOrdersAction   `json:"cxa,omitempty"`
	StopOrder          *StopOrderAction         `json:"st,omitempty"`
	TakeProfitOrder    *TakeProfitAction        `json:"tp,omitempty"`
	ApproveBuilderCode *BuilderCodeAction       `json:"abc,omitempty"`
	RevokeBuilderCode  *RevokeBuilderCodeAction `json:"rbc,omitempty"`
	Transfer           *TransferAction          `json:"transfer,omitempty"`
}

type LimitOrderAction struct {
	Price           float64            `json:"px"`
	TimeInForce     TimeInForce        `json:"tif"`
	Symbol          string             `json:"c"`
	Buy             bool               `json:"b"`
	Size            float64            `json:"sz"`
	ReduceOnly      bool               `json:"r"`
	IsolatedAccount bool               `json:"i"`
	BuilderCode     *BuilderCodeAction `json:"builderCode,omitempty"`
}

type TimeInForce uint32

const (
	GTC TimeInForce = iota
	IOC
	ALO
)

func (LimitOrderAction) Discriminant() uint32 {
	return actionLimit
}

type MarketOrderAction struct {
	Symbol          string             `json:"c"`
	Buy             bool               `json:"b"`
	Size            float64            `json:"sz"`
	ReduceOnly      bool               `json:"r"`
	IsolatedAccount bool               `json:"i"`
	BuilderCode     *BuilderCodeAction `json:"builderCode,omitempty"`
}

func (MarketOrderAction) Discriminant() uint32 {
	return actionMarket
}

type ModifyOrderAction struct {
	OrderID string  `json:"oid"` //the base58 hash
	Symbol  string  `json:"symbol"`
	Amount  float64 `json:"amount"`
}

func (ModifyOrderAction) Discriminant() uint32 {
	return actionModify
}

type CancelSingleOrderAction struct {
	Symbol  string `json:"c"`
	OrderID string `json:"oid"`
}

func (CancelSingleOrderAction) Discriminant() uint32 {
	return actionCancel
}

func (c CancelSingleOrderAction) GetSymbol() string {
	return c.Symbol
}

type CancelAllOrdersAction struct {
	Symbols []string `json:"c"`
}

func (CancelAllOrdersAction) Discriminant() uint32 {
	return actionCancelAll
}

type StopOrderAction struct {
	Symbol           string   `json:"c"`
	TriggerDirection bool     `json:"d"`
	Size             float64  `json:"sz"`
	TriggerPrice     float64  `json:"tr"`
	LimitPrice       *float64 `json:"lim,omitempty"`
	IsolatedAccount  bool     `json:"i"`
}

func (StopOrderAction) Discriminant() uint32 {
	return actionStop
}

type TakeProfitAction struct {
	StopOrderAction
}

func (TakeProfitAction) Discriminant() uint32 {
	return actionTakeProfit
}

type BuilderCodeAction struct {
	To  string `json:"to"`
	Fee int    `json:"fee"`
}

func (BuilderCodeAction) Discriminant() uint32 {
	return actionApproveBuilderCode
}

type RevokeBuilderCodeAction struct {
	To string `json:"to"`
}

func (RevokeBuilderCodeAction) Discriminant() uint32 {
	return actionRevokeBuilderCode
}

type TransferKind string

const (
	Internal TransferKind = "internal"
	External TransferKind = "external"
)

type TransferAction struct {
	Kind TransferKind `json:"k"`
	From string       `json:"from"`
	To   string       `json:"to"`
	// eg USDC
	MarginSymbol string  `json:"marginSymbol"`
	MarginAmount float64 `json:"marginAmount"`
}

func (TransferAction) Discriminant() uint32 {
	return actionTransfer
}

// struct to send to bulk
type SignMessage struct {
	Actions   []Action `json:"actions"`
	Nonce     uint64   `json:"nonce"`
	Account   string   `json:"account"`
	Signer    string   `json:"signer"`
	Signature string   `json:"signature"`
}

type SignInput struct {
	Actions []Action `json:"actions"`
	Nonce   uint64   `json:"nonce"`
	Account string   `json:"account"`
}

type DomainByte uint8

const (
	Mainnet DomainByte = iota + 1
	Testnet
	Devnet
)

type KeyPair struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

type Signer struct {
	keypair *KeyPair
	domain  DomainByte
}

const (
	actionMarket     uint32 = iota
	actionLimit             // 1
	actionModify            // 2
	actionCancel            // 3
	actionCancelAll         // 4
	actionStop              // 5
	actionTakeProfit        // 6
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	actionTransfer //29
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	actionApproveBuilderCode //40
	actionRevokeBuilderCode  //41
)
