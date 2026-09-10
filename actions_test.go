package bulkkeychain

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

func TestWriteMarketOrderWithoutBuilderCodeAction(t *testing.T) {
	action := &MarketOrderAction{
		Symbol:          "BTC-USD",
		Buy:             true,
		Size:            0.1,
		ReduceOnly:      false,
		IsolatedAccount: true,
		BuilderCode:     nil,
	}

	var got bytes.Buffer
	err := writeAction(&got, action)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x00, 0x00, 0x00, 0x00, //market order discriminator

		// symbol length: uint64(7) // B T C - U S D
		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		'B', 'T', 'C', '-', 'U', 'S', 'D',

		// buy=true
		0x01,

		// 0.1 * 1e8 = 10_000_000 as uint64 little endian
		0x80, 0x96, 0x98, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x00, // reduce only = false

		0x01, // isolated  account= true
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("wrong serialization\n i got %x, when i want %x", got.Bytes(), want)
	}

}

func TestWriteMarketOrderWithBuilderCodeAction(t *testing.T) {
	recpBytes := []byte("gurtcryptogurtcryptogurtcryptogu")
	address := base58.Encode([]byte(recpBytes)) // 32 length

	action := &MarketOrderAction{
		Symbol:          "BTC-USD",
		Buy:             true,
		Size:            0.1,
		ReduceOnly:      false,
		IsolatedAccount: true,
		BuilderCode: &BuilderCodeAction{
			To:  address,
			Fee: 10,
		},
	}

	var got bytes.Buffer
	err := writeAction(&got, action)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x00, 0x00, 0x00, 0x00, //market order discriminator

		// symbol length: uint64(7) // B T C - U S D
		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		'B', 'T', 'C', '-', 'U', 'S', 'D',

		// buy=true
		0x01,

		// 0.1 * 1e8 = 10_000_000 as uint64 little endian
		0x80, 0x96, 0x98, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x00, // reduce only = false

		0x01, // isolated  account= true
	}

	want = append(want, 0x01)
	want = append(want, recpBytes...)
	want = append(want, 0x0a) // fee

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("wrong serialization\n i got %x,\ni want %x", got.Bytes(), want)
	}

}

func TestLimitOrderAction(t *testing.T) {
	action := &LimitOrderAction{
		Symbol:          "BTC-USD",
		Buy:             true,
		Size:            0.1,
		ReduceOnly:      false,
		IsolatedAccount: true,
		BuilderCode:     nil,

		Price:       100_000,
		TimeInForce: IOC,
	}

	var got bytes.Buffer

	err := writeAction(&got, action)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x01, 0x00, 0x00, 0x00,

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		'B', 'T', 'C', '-', 'U', 'S', 'D',

		0x01,

		0x00, 0xA0, 0x72, 0x4E,
		0x18, 0x09, 0x00, 0x00,

		0x80, 0x96, 0x98, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x01, 0x00, 0x00, 0x00,

		0x00,
		0x01,
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("wrong serialization\n expected: %x,\n got %x", want, got.Bytes())
	}
}

func TestModifyOrderAction(t *testing.T) {
	recpBytes := []byte("gurtcryptogurtcryptogurtcryptogu")
	orderId := base58.Encode([]byte(recpBytes)) // 32 length

	action := &ModifyOrderAction{
		OrderID: orderId,
		Symbol:  "BTC-USD",
		Amount:  0.1,
	}
	var got bytes.Buffer

	err := writeAction(&got, action)

	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x02, 0x00, 0x00, 0x00,
	}
	want = append(want, recpBytes...)

	symbol := []byte{
		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	want = append(want, symbol...)
	want = append(want, []byte("BTC-USD")...)
	amount := []byte{
		// plain float64(0.1)
		0x9a, 0x99, 0x99, 0x99,
		0x99, 0x99, 0xb9, 0x3f,
	}

	want = append(want, amount...)

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestCancelOrderAction(t *testing.T) {
	recpBytes := []byte("gurtcryptogurtcryptogurtcryptogu")
	orderId := base58.Encode([]byte(recpBytes)) // 32 length
	action := &CancelSingleOrderAction{
		Symbol:  "BTC-USD",
		OrderID: orderId,
	}

	var got bytes.Buffer
	err := writeAction(&got, action)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{
		0x03, 0x00, 0x00, 0x00,

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		'B', 'T', 'C', '-', 'U', 'S', 'D',
	}

	expected = append(expected, recpBytes...)

	if !bytes.Equal(expected, got.Bytes()) {
		t.Fatalf("expected: %x,\n got %x", expected, got.Bytes())
	}
}

func TestCancelAllOrdersAction(t *testing.T) {
	action := &CancelAllOrdersAction{
		Symbols: []string{"BTC-USD", "ETH-USD"},
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x04, 0x00, 0x00, 0x00,

		0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		'B', 'T', 'C', '-', 'U', 'S', 'D',

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		'E', 'T', 'H', '-', 'U', 'S', 'D',
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestStopOrderAction(t *testing.T) {
	action := &StopOrderAction{
		Symbol:           "BTC-USD",
		TriggerDirection: false,
		Size:             0.1,
		TriggerPrice:     100_000,
		LimitPrice:       nil,
		IsolatedAccount:  false,
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x05, 0x00, 0x00, 0x00,

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		'B', 'T', 'C', '-', 'U', 'S', 'D',

		0x00,

		0x80, 0x96, 0x98, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x00, 0xa0, 0x72, 0x4e,
		0x18, 0x09, 0x00, 0x00,

		0x00,
		0x00,
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestTakeProfitOrderAction(t *testing.T) {
	limitPrice := 100_000.0
	action := &TakeProfitAction{
		StopOrderAction: StopOrderAction{
			Symbol:           "BTC-USD",
			TriggerDirection: true,
			Size:             0.1,
			TriggerPrice:     100_000,
			LimitPrice:       &limitPrice,
			IsolatedAccount:  true,
		},
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x06, 0x00, 0x00, 0x00,

		0x07, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		'B', 'T', 'C', '-', 'U', 'S', 'D',

		0x01,

		0x80, 0x96, 0x98, 0x00,
		0x00, 0x00, 0x00, 0x00,

		0x00, 0xa0, 0x72, 0x4e,
		0x18, 0x09, 0x00, 0x00,

		0x01,
		0x00, 0xa0, 0x72, 0x4e,
		0x18, 0x09, 0x00, 0x00,

		0x01,
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestApproveBuilderCodeAction(t *testing.T) {
	recipient := []byte("gurtcryptogurtcryptogurtcryptogu")
	action := &BuilderCodeAction{
		To:  base58.Encode(recipient),
		Fee: 10,
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{0x28, 0x00, 0x00, 0x00}
	want = append(want, recipient...)
	want = append(want, 0x0a)

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestRevokeBuilderCodeAction(t *testing.T) {
	recipient := []byte("gurtcryptogurtcryptogurtcryptogu")
	action := &RevokeBuilderCodeAction{
		To: base58.Encode(recipient),
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{0x29, 0x00, 0x00, 0x00}
	want = append(want, recipient...)

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestTransferAction(t *testing.T) {
	from := []byte("gurtcryptogurtcryptogurtcryptogu")
	to := []byte("bulkcryptobulkcryptobulkcryptobu")
	action := &TransferAction{
		Kind:         External,
		From:         base58.Encode(from),
		To:           base58.Encode(to),
		MarginSymbol: "USDC",
		MarginAmount: 100,
	}

	var got bytes.Buffer
	if err := writeAction(&got, action); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x1d, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
	}
	want = append(want, from...)
	want = append(want, to...)
	want = append(want,
		0x04, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		'U', 'S', 'D', 'C',
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x59, 0x40,
	)

	if !bytes.Equal(want, got.Bytes()) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}

func TestAgentWalletCreation(t *testing.T) {
	pubKy := []byte("gurtcryptogurtcryptogurtcryptogu")

	action := &AgentWalletCreationAction{
		PublicKey: base58.Encode(pubKy),
		Delete:    false,
	}

	var got bytes.Buffer
	err := writeAgentWalletAction(&got, action)
	if err != nil {
		t.Fatalf("unable to write action: %v", err)
	}

	want := []byte{}
	want = append(want, pubKy...)
	want = append(want, 0x00)

	if !bytes.Equal(got.Bytes(), want) {
		t.Fatalf("expected: %x,\n got: %x", want, got.Bytes())
	}
}
