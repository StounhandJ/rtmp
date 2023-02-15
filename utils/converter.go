package utils

import "encoding/binary"

func Uint32ToByteArray(val uint32, size int) []byte {
	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf[:], val)

	return buf[:size]
}
