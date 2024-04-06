package modbus

import (
	"encoding/binary"
)

func convertTooLargeNumber(in []byte) (uint32) {
	var u32 uint32
	var out []uint32

	for i := 0; i < len(in); i++ {
		in[i] = ^in[i]
	}

	for i := 0; i < len(in); i += 4 {
		u32 = binary.BigEndian.Uint32(in[i : i+4])

		out = append(out, u32)
	}

	return out[0]
}
