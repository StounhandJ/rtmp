package Handshake

import (
	"encoding/binary"
	"io"
	"time"
)

type S2C2 struct {
	Time   uint32
	Time2  uint32
	Random [1528]byte
}

func NewS2C2(timeIn uint32, random [1528]byte) *S2C2 {
	sc := &S2C2{
		Time:  timeIn,
		Time2: uint32(time.Now().UnixNano() / int64(time.Millisecond)),
	}
	copy(sc.Random[:], random[:])
	return sc
}

func DecodeS2C2(reader io.Reader) (*S2C2, error) {
	sc := &S2C2{}

	var buf [4]byte

	if _, err := io.ReadAtLeast(reader, buf[:], len(buf)); err != nil {
		return nil, err
	}
	sc.Time = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(reader, buf[:], len(buf)); err != nil {
		return nil, err
	}
	sc.Time2 = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(reader, sc.Random[:], len(sc.Random)); err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *S2C2) EncodeS2C2(writer io.Writer) error {
	buf := [4]byte{}

	binary.BigEndian.PutUint32(buf[:], sc.Time)
	if _, err := writer.Write(buf[:]); err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf[:], sc.Time2)
	if _, err := writer.Write(buf[:]); err != nil {
		return err
	}

	if _, err := writer.Write(sc.Random[:]); err != nil {
		return err
	}

	return nil
}
