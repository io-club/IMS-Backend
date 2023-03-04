package util

import (
	"crypto/rand"
	"encoding/binary"
)

func RandomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}
