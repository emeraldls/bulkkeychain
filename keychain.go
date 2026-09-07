package bulkkeychain

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

// https://docs.bulk.trade/api-reference/signing#what-gets-signed
func SignOrder(orderInput OrderInput, privateKey ed25519.PrivateKey, domain DomainByte) (Order, error) {
	actions, err := serializeActions(orderInput.Actions)
	if err != nil {
		return Order{}, err
	}

	account := base58.Decode(orderInput.Account)

	var message = bytes.Buffer{}
	message.Write(actions)

	err = binary.Write(&message, binary.LittleEndian, orderInput.Nonce)
	if err != nil {
		return Order{}, err
	}

	message.Write(account)
	message.WriteByte(byte(domain))

	signature := ed25519.Sign(privateKey, message.Bytes())
	signer := privateKey.Public().(ed25519.PublicKey)

	var order Order
	order.Actions = orderInput.Actions
	order.Account = orderInput.Account
	order.Nonce = orderInput.Nonce
	order.Signer = base58.Encode(signer)
	order.Signature = base58.Encode(signature)

	return order, nil
}

func serializeActions(actions []Action) ([]byte, error) {
	if len(actions) == 0 {
		return nil, errors.New("add at least one action")
	}

	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, uint64(len(actions)))
	if err != nil {
		return nil, fmt.Errorf("unablet to write actions: %w", err)
	}

	for _, action := range actions {
		orderAction := unwrapAction(action)
		err = writeAction(&buf, orderAction)
		if err != nil {
			return nil, fmt.Errorf("action: %T, err: %w", orderAction, err)
		}
	}

	return buf.Bytes(), nil
}

func unwrapAction(action Action) OrderAction {
	var selectedAction OrderAction

	if action.MarketOrder != nil {
		selectedAction = action.MarketOrder
	}
	if action.LimitOrder != nil {
		selectedAction = action.LimitOrder
	}
	if action.ModifyOrder != nil {
		selectedAction = action.ModifyOrder
	}
	if action.CancelOrder != nil {
		selectedAction = action.CancelOrder
	}
	if action.CancelAllOrders != nil {
		selectedAction = action.CancelAllOrders
	}
	if action.StopOrder != nil {
		selectedAction = action.StopOrder
	}
	if action.TakeProfitOrder != nil {
		selectedAction = action.TakeProfitOrder
	}

	return selectedAction
}
