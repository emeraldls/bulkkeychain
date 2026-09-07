package bulkkeychain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func writeAction(buf *bytes.Buffer, signAction Signing) error {
	if err := binary.Write(buf, binary.LittleEndian, signAction.Discriminant()); err != nil {
		return fmt.Errorf("unable to write action discriminant: %w", err)
	}

	switch v := signAction.(type) {
	case *MarketOrderAction:
		return writeMarketOrder(buf, v)

	case *LimitOrderAction:
		return writeLimitOrder(buf, v)

	case *ModifyOrderAction:
		return writeModifyOrder(buf, v)

	case *CancelSingleOrderAction:
		return writeCancelSingleOrder(buf, v)

	case *CancelAllOrdersAction:
		return writeCancelAllOrders(buf, v)

	case *StopOrderAction:
		return writeStopOrder(buf, v)

	case *TakeProfitAction:
		return writeTakeProfitOrder(buf, v)

	case *BuilderCodeAction:
		return writeApproveBuilderCodeAction(buf, v)

	case *RevokeBuilderCodeAction:
		return writeRevokeBuilderCodeAction(buf, v)

	case *TransferAction:
		return writeTransferAction(buf, v)

	default:
		return fmt.Errorf("unsupported action type: %T", v)
	}
}

// https://docs.bulk.trade/api-reference/signing#marketorder-discriminant-0
func writeMarketOrder(buf *bytes.Buffer, order *MarketOrderAction) error {
	err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbol)))
	if err != nil {
		return fmt.Errorf("unable to write symbol: %s", err.Error())
	}

	buf.WriteString(order.Symbol)

	writeBool(buf, order.Buy)
	err = writeScaledFloat(buf, order.Size)
	if err != nil {
		return fmt.Errorf("size: %v", err)
	}

	writeBool(buf, order.ReduceOnly)
	writeBool(buf, order.IsolatedAccount)

	if order.BuilderCode == nil {
		return nil
	}

	recipient := base58.Decode(order.BuilderCode.To)
	if len(recipient) != 32 {
		return fmt.Errorf("builder code must have a length of 32bbyte: current length: %d", len(recipient))
	}

	buf.WriteByte(1)
	buf.Write(recipient)
	buf.WriteByte(byte(order.BuilderCode.Fee))

	return nil
}

// https://docs.bulk.trade/api-reference/signing#limitorder-discriminant-1
func writeLimitOrder(buf *bytes.Buffer, order *LimitOrderAction) error {
	err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbol)))
	if err != nil {
		return fmt.Errorf("unable to write symbol: %s", err.Error())
	}

	buf.WriteString(order.Symbol)

	writeBool(buf, order.Buy)
	err = writeScaledFloat(buf, order.Price)
	if err != nil {
		return fmt.Errorf("price: %v", err)
	}

	err = writeScaledFloat(buf, order.Size)
	if err != nil {
		return fmt.Errorf("size: %v", err)
	}

	err = binary.Write(buf, binary.LittleEndian, uint32(order.TimeInForce))
	if err != nil {
		return fmt.Errorf("unable to write tif: %w", err)
	}

	writeBool(buf, order.ReduceOnly)
	writeBool(buf, order.IsolatedAccount)

	if order.BuilderCode == nil {
		return nil
	}

	recipient := base58.Decode(order.BuilderCode.To)
	if len(recipient) != 32 {
		return fmt.Errorf("builder code must have a length of 32bbyte: current length: %d", len(recipient))
	}

	buf.WriteByte(1)
	buf.Write(recipient)
	buf.WriteByte(byte(order.BuilderCode.Fee))

	return nil
}

// https://docs.bulk.trade/api-reference/signing#modifyorder-discriminant-2
func writeModifyOrder(buf *bytes.Buffer, order *ModifyOrderAction) error {
	orderId := base58.Decode(order.OrderID)
	if len(orderId) != 32 {
		return fmt.Errorf("invalid order Id, expected 32 bytes, got: %d", len(orderId))
	}

	buf.Write(orderId)
	err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbol)))
	if err != nil {
		return fmt.Errorf("unable to write symbol: %w", err)
	}

	buf.WriteString(order.Symbol)
	err = binary.Write(buf, binary.LittleEndian, order.Amount)
	if err != nil {
		return fmt.Errorf("unable tow write moidified amount: %w", err)
	}

	return nil
}

// https://docs.bulk.trade/api-reference/signing#cancel-discriminant-3
func writeCancelSingleOrder(buf *bytes.Buffer, order *CancelSingleOrderAction) error {
	err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbol)))
	if err != nil {
		return fmt.Errorf("unable to write symbol: %w", err)
	}

	buf.WriteString(order.Symbol)

	orderId := base58.Decode(order.OrderID)
	if len(orderId) != 32 {
		return fmt.Errorf("invalid order Id, expected 32byte, got %d", len(orderId))
	}

	buf.Write(orderId)
	return nil
}

// https://docs.bulk.trade/api-reference/signing#cancelall-discriminant-4
func writeCancelAllOrders(buf *bytes.Buffer, order *CancelAllOrdersAction) error {
	if err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbols))); err != nil {
		return fmt.Errorf("unable to write symbols count: %w", err)
	}

	for _, symbol := range order.Symbols {
		if err := binary.Write(buf, binary.LittleEndian, uint64(len(symbol))); err != nil {
			return fmt.Errorf("unable to write symbol: %w", err)
		}

		buf.WriteString(symbol)
	}

	return nil
}

// https://docs.bulk.trade/api-reference/signing#stop-discriminant-5-and-takeprofit-discriminant-6
func writeStopOrder(buf *bytes.Buffer, order *StopOrderAction) error {
	if err := binary.Write(buf, binary.LittleEndian, uint64(len(order.Symbol))); err != nil {
		return fmt.Errorf("unable to write symbol: %w", err)
	}

	buf.WriteString(order.Symbol)
	writeBool(buf, order.TriggerDirection)

	if err := writeScaledFloat(buf, order.Size); err != nil {
		return fmt.Errorf("size: %w", err)
	}

	if err := writeScaledFloat(buf, order.TriggerPrice); err != nil {
		return fmt.Errorf("trigger price: %w", err)
	}

	if order.LimitPrice == nil {
		buf.WriteByte(0)
	} else {
		buf.WriteByte(1)
		if err := writeScaledFloat(buf, *order.LimitPrice); err != nil {
			return fmt.Errorf("limit price: %w", err)
		}
	}

	writeBool(buf, order.IsolatedAccount)

	return nil
}

// https://docs.bulk.trade/api-reference/signing#stop-discriminant-5-and-takeprofit-discriminant-6
func writeTakeProfitOrder(buf *bytes.Buffer, order *TakeProfitAction) error {
	return writeStopOrder(buf, &order.StopOrderAction)
}

// https://docs.bulk.trade/api-reference/signing#approvebuildercode-discriminant-40
func writeApproveBuilderCodeAction(buf *bytes.Buffer, action *BuilderCodeAction) error {
	recp := base58.Decode(action.To)
	if len(recp) != 32 {
		return errors.New("builder recipient length should be 32 bytes")
	}

	buf.Write(recp)
	buf.WriteByte(byte(action.Fee))

	return nil
}

// https://docs.bulk.trade/api-reference/signing#revokebuildercode-discriminant-41
func writeRevokeBuilderCodeAction(buf *bytes.Buffer, action *RevokeBuilderCodeAction) error {
	recp := base58.Decode(action.To)
	if len(recp) != 32 {
		return errors.New("builder recipient length should be 32 bytes")
	}

	buf.Write(recp)
	return nil
}

// https://docs.bulk.trade/api-reference/signing#transfer-discriminant-29
func writeTransferAction(buf *bytes.Buffer, action *TransferAction) error {
	var kind uint32
	switch action.Kind {
	case Internal:
		kind = 0
	case External:
		kind = 1
	default:
		return fmt.Errorf("invalid transfer kind: %q", action.Kind)
	}

	from := base58.Decode(action.From)
	if len(from) != 32 {
		return fmt.Errorf("from pubkey must be 32 bytes, got %d", len(from))
	}

	to := base58.Decode(action.To)
	if len(to) != 32 {
		return fmt.Errorf("to pubkey must be 32 bytes, got %d", len(to))
	}

	err := binary.Write(buf, binary.LittleEndian, kind)
	if err != nil {
		return fmt.Errorf("unable to write transfer kind: %w", err)
	}

	buf.Write(from)
	buf.Write(to)

	err = binary.Write(buf, binary.LittleEndian, uint64(len(action.MarginSymbol)))
	if err != nil {
		return fmt.Errorf("unable to write margin symbol: %w", err)
	}
	buf.WriteString(action.MarginSymbol)

	err = binary.Write(buf, binary.LittleEndian, action.MarginAmount)
	if err != nil {
		return fmt.Errorf("unable to write margin amount: %w", err)
	}

	return nil
}
