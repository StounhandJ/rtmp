package Handshake

import (
	"encoding/binary"
	"io"
	"math/rand"
	"time"
)

type S1C1 struct {
	Time   uint32
	Zero   [4]byte
	Random [1528]byte
}

func NewS1C1() *S1C1 {
	sc := &S1C1{
		Time: uint32(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	if _, err := rand.Read(sc.Random[:]); err != nil { // Random Seq
		return nil
	}

	return sc
}

func DecodeS1C1(reader io.Reader) (*S1C1, error) {
	sc := &S1C1{}

	var buf [4]byte

	if _, err := io.ReadAtLeast(reader, buf[:], len(buf)); err != nil {
		return nil, err
	}
	sc.Time = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(reader, sc.Zero[:], len(sc.Zero)); err != nil {
		return nil, err
	}

	if _, err := io.ReadAtLeast(reader, sc.Random[:], len(sc.Random)); err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *S1C1) EncodeS1C1(writer io.Writer) error {
	buf := [4]byte{}

	binary.BigEndian.PutUint32(buf[:], sc.Time)
	if _, err := writer.Write(buf[:]); err != nil {
		return err
	}

	if _, err := writer.Write(sc.Zero[:]); err != nil {
		return err
	}

	if _, err := writer.Write(sc.Random[:]); err != nil {
		return err
	}

	return nil
}
