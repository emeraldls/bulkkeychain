package bulkkeychain

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func NewSigner(keypair *KeyPair, domain DomainByte) *Signer {
	return &Signer{
		keypair,
		domain,
	}
}

// https://docs.bulk.trade/api-reference/signing#what-gets-signed
func (s *Signer) Sign(signInput SignInput) (SignMessage, error) {
	if len(s.keypair.privateKey) != ed25519.PrivateKeySize {
		return SignMessage{}, fmt.Errorf("private key must be %d bytes", ed25519.PrivateKeySize)
	}

	actions, err := serializeActions(signInput.Actions)
	if err != nil {
		return SignMessage{}, err
	}

	account := base58.Decode(signInput.Account)

	var message = bytes.Buffer{}
	message.Write(actions)

	err = binary.Write(&message, binary.LittleEndian, signInput.Nonce)
	if err != nil {
		return SignMessage{}, err
	}

	message.Write(account)
	message.WriteByte(byte(s.domain))

	signature := ed25519.Sign(s.keypair.privateKey, message.Bytes())
	signer := s.keypair.privateKey.Public().(ed25519.PublicKey)

	var order SignMessage
	order.Actions = signInput.Actions
	order.Account = signInput.Account
	order.Nonce = signInput.Nonce
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
		orderAction, err := unwrapAction(action)
		if err != nil {
			return nil, fmt.Errorf("action %w", err)
		}

		err = writeAction(&buf, orderAction)
		if err != nil {
			return nil, fmt.Errorf("action: %T, err: %w", orderAction, err)
		}
	}

	return buf.Bytes(), nil
}

func unwrapAction(action Action) (Signing, error) {
	actions := make([]Signing, 0, 1)

	if action.MarketOrder != nil {
		actions = append(actions, action.MarketOrder)
	}
	if action.LimitOrder != nil {
		actions = append(actions, action.LimitOrder)
	}
	if action.ModifyOrder != nil {
		actions = append(actions, action.ModifyOrder)
	}
	if action.CancelOrder != nil {
		actions = append(actions, action.CancelOrder)
	}
	if action.CancelAllOrders != nil {
		actions = append(actions, action.CancelAllOrders)
	}
	if action.StopOrder != nil {
		actions = append(actions, action.StopOrder)
	}
	if action.TakeProfitOrder != nil {
		actions = append(actions, action.TakeProfitOrder)
	}
	if action.ApproveBuilderCode != nil {
		actions = append(actions, action.ApproveBuilderCode)
	}
	if action.RevokeBuilderCode != nil {
		actions = append(actions, action.RevokeBuilderCode)
	}
	if action.Transfer != nil {
		actions = append(actions, action.Transfer)
	}

	if action.AgentWallet != nil {
		actions = append(actions, action.AgentWallet)
	}

	if len(actions) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one action type, got %d",
			len(actions),
		)
	}

	return actions[0], nil
}
