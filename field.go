package aof

import (
	"encoding/binary"
)

const (
	SizeUint16 = 2
	SizeUint24 = 3
	SizeUint32 = 4
	SizeUint48 = 6
	SizeUint64 = 8
)

func DecodeUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func DecodeUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func DecodeUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
