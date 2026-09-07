package bulkkeychain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func writeAction(buf *bytes.Buffer, orderAction OrderAction) error {
	if err := binary.Write(buf, binary.LittleEndian, orderAction.Discriminant()); err != nil {
		return fmt.Errorf("unable to write action discriminant: %w", err)
	}

	switch v := orderAction.(type) {
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

	default:
		return fmt.Errorf("unsupported order action type: %T", v)
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
