package bulkkeychain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

// no need for branching
func writeBool(buf *bytes.Buffer, value bool) {
	buf.WriteByte(*(*byte)(unsafe.Pointer(&value)))

	/*
		alternative for branching
		if value {
			buf.WriteByte(1)
			return
		}

		buf.WriteByte(0)
	*/
}

func writeScaledFloat(buf *bytes.Buffer, value float64) error {
	scaled := math.Round(value * 1e8)

	if math.IsNaN(scaled) || math.IsInf(scaled, 0) || scaled < 0 || scaled >= math.Exp2(64) {
		return fmt.Errorf("value (%f) cannot be scaled as u64", value)
	}

	return binary.Write(buf, binary.LittleEndian, uint64(scaled))
}
