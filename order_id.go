package bulkkeychain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func ComputeOrderID(action Action, account string, nonce uint64, index uint32) (string, error) {
	order, err := unwrapAction(action)
	if err != nil {
		return "", err
	}

	switch v := order.(type) {
	case *LimitOrderAction:
		copy := *v
		copy.BuilderCode = nil
		order = &copy
	case *MarketOrderAction:
		copy := *v
		copy.BuilderCode = nil
		order = &copy
	case *StopOrderAction:
		copy := *v
		copy.BuilderCode = nil
		order = &copy
	case *TakeProfitAction:
		copy := *v
		copy.BuilderCode = nil
		order = &copy
	default:
		return "", fmt.Errorf("action does not create an order: %T", order)
	}
	owner := base58.Decode(account)
	var encoded bytes.Buffer
	if err := writeAction(&encoded, order); err != nil {
		return "", err
	}
	preimage := binary.LittleEndian.AppendUint32(nil, index)
	preimage = append(preimage, encoded.Bytes()...)
	preimage = append(preimage, owner...)
	preimage = binary.LittleEndian.AppendUint64(preimage, nonce)
	hash := sha256.Sum256(preimage)
	return base58.Encode(hash[:]), nil
}
