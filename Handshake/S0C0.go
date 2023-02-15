package Handshake

import "io"

type S0C0 struct {
	Version byte
}

func NewS0C0(version byte) *S0C0 {
	return &S0C0{
		Version: version,
	}
}

func DecodeS0C0(reader io.Reader) (*S0C0, error) {
	sc := &S0C0{}

	buf := [1]byte{}

	if _, err := io.ReadAtLeast(reader, buf[:], 1); err != nil {
		return nil, err
	}
	sc.Version = buf[0]

	return sc, nil
}

func (sc *S0C0) EncodeS0C0(writer io.Writer) error {
	buf := [1]byte{sc.Version}

	_, err := writer.Write(buf[:])
	if err != nil {
		return err
	}

	return nil
}
