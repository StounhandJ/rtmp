package utils

import "encoding/binary"

func Uint32ToByteArray(val uint32, size int) []byte {
	buf := make([]byte, 4)

	if size >= 4 {
		binary.LittleEndian.PutUint32(buf[:], val)

	} else {
		binary.BigEndian.PutUint32(buf[:], val)
	}

	return buf[4-size:]
}
