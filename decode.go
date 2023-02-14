package main

import (
	"encoding/binary"
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (d *Decoder) DecodeS0C0(h *S0C0) error {
	buf := [1]byte{}

	if _, err := io.ReadAtLeast(d.r, buf[:], 1); err != nil {
		return err
	}
	*h = S0C0(buf[0])

	return nil
}

func (d *Decoder) DecodeS1C1(h *S1C1) error {
	var buf [4]byte

	if _, err := io.ReadAtLeast(d.r, buf[:], len(buf)); err != nil {
		return err
	}
	h.Time = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(d.r, h.Version[:], len(h.Version)); err != nil {
		return err
	}

	if _, err := io.ReadAtLeast(d.r, h.Random[:], len(h.Random)); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) DecodeS2C2(h *S2C2) error {
	var buf [4]byte

	if _, err := io.ReadAtLeast(d.r, buf[:], len(buf)); err != nil {
		return err
	}
	h.Time = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(d.r, buf[:], len(buf)); err != nil {
		return err
	}
	h.Time2 = binary.BigEndian.Uint32(buf[:])

	if _, err := io.ReadAtLeast(d.r, h.Random[:], len(h.Random)); err != nil {
		return err
	}

	return nil
}

//  ----------------CHUNK----------------

type chunkBasicHeader struct {
	fmt           byte
	chunkStreamID int /* [0, 65599] */
}

func decodeChunkBasicHeader(r io.Reader, buf []byte, bh *chunkBasicHeader) error {
	if buf == nil || len(buf) < 3 {
		buf = make([]byte, 3)
	}

	if _, err := io.ReadAtLeast(r, buf[:1], 1); err != nil {
		return err
	}

	fmtTy := (buf[0] & 0xc0) >> 6 // 0b11000000 >> 6
	csID := int(buf[0] & 0x3f)    // 0b00111111

	switch csID {
	case 0:
		if _, err := io.ReadAtLeast(r, buf[1:2], 1); err != nil {
			return err
		}
		csID = int(buf[1]) + 64

	case 1:
		if _, err := io.ReadAtLeast(r, buf[1:], 2); err != nil {
			return err
		}
		csID = int(buf[2])*256 + int(buf[1]) + 64
	}

	bh.fmt = fmtTy
	bh.chunkStreamID = csID

	return nil
}

type chunkMessageHeader struct {
	timestamp       uint32 // fmt = 0
	timestampDelta  uint32 // fmt = 1 | 2
	messageLength   uint32 // fmt = 0 | 1
	messageTypeID   byte   // fmt = 0 | 1
	messageStreamID uint32 // fmt = 0
}

func decodeChunkMessageHeader(r io.Reader, fmt byte, buf []byte, mh *chunkMessageHeader) error {
	if buf == nil || len(buf) < 11 {
		buf = make([]byte, 11)
	}
	cache32bits := make([]byte, 4)

	switch fmt {
	case 0:
		if _, err := io.ReadAtLeast(r, buf[:11], 11); err != nil {
			return err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.timestamp = binary.BigEndian.Uint32(cache32bits)
		copy(cache32bits[1:], buf[3:6]) // 24bits BE
		mh.messageLength = binary.BigEndian.Uint32(cache32bits)
		mh.messageTypeID = buf[6]                                  // 8bits
		mh.messageStreamID = binary.LittleEndian.Uint32(buf[7:11]) // 32bits

		if mh.timestamp == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return err
			}
			mh.timestamp = binary.BigEndian.Uint32(cache32bits)
		}

	case 1:
		if _, err := io.ReadAtLeast(r, buf[:7], 7); err != nil {
			return err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.timestampDelta = binary.BigEndian.Uint32(cache32bits)
		copy(cache32bits[1:], buf[3:6]) // 24bits BE
		mh.messageLength = binary.BigEndian.Uint32(cache32bits)
		mh.messageTypeID = buf[6] // 8bits

		if mh.timestampDelta == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return err
			}
			mh.timestampDelta = binary.BigEndian.Uint32(cache32bits)
		}

	case 2:
		if _, err := io.ReadAtLeast(r, buf[:3], 3); err != nil {
			return err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.timestampDelta = binary.BigEndian.Uint32(cache32bits)

		if mh.timestampDelta == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return err
			}
			mh.timestampDelta = binary.BigEndian.Uint32(cache32bits)
		}

	case 3:
		// DO NOTHING

	default:
		panic("Unexpected fmt")
	}

	return nil
}
