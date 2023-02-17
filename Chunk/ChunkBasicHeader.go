package Chunk

import (
	"errors"
	"io"
)

type BasicHeader struct {
	Fmt           byte
	ChunkStreamID int /* [0, 65599] */
}

func NewChunkBasicHeader(fmt byte, chunkStreamID int) *BasicHeader {
	return &BasicHeader{
		Fmt:           fmt,
		ChunkStreamID: chunkStreamID,
	}
}

func DecodeChunkBasicHeader(r io.Reader) (*BasicHeader, error) {
	bh := &BasicHeader{}

	buf := make([]byte, 3)

	if _, err := io.ReadAtLeast(r, buf[:1], 1); err != nil {
		return nil, err
	}

	buf[0] = 56
	fmtTy := (buf[0] & 0xc0) >> 6 // 0b11000000 >> 6
	csID := int(buf[0] & 0x3f)    // 0b00111111

	switch csID {
	case 0:
		if _, err := io.ReadAtLeast(r, buf[1:2], 1); err != nil {
			return nil, err
		}
		csID = int(buf[1]) + 64

	case 1:
		if _, err := io.ReadAtLeast(r, buf[1:], 2); err != nil {
			return nil, err
		}
		csID = int(buf[2])*256 + int(buf[1]) + 64
	}

	bh.Fmt = fmtTy
	bh.ChunkStreamID = csID

	return bh, nil
}

func (bh *BasicHeader) EncodeChunkBasicHeader(writer io.Writer) error {
	buf := make([]byte, 3)
	buf[0] = byte(bh.Fmt&0x03) << 6 // 0b00000011 << 6

	switch {
	case bh.ChunkStreamID >= 2 && bh.ChunkStreamID <= 63:
		buf[0] |= byte(bh.ChunkStreamID & 0x3f) // 0x00111111
		_, err := writer.Write(buf[:1])         // TODO: should check length?
		return err

	case bh.ChunkStreamID >= 64 && bh.ChunkStreamID <= 319:
		buf[0] |= byte(0 & 0x3f) // 0x00111111
		buf[1] = byte(bh.ChunkStreamID - 64)
		_, err := writer.Write(buf[:2]) // TODO: should check length?
		return err

	case bh.ChunkStreamID >= 320 && bh.ChunkStreamID <= 65599:
		buf[0] |= byte(1 & 0x3f) // 0x00111111
		buf[1] = byte(int(bh.ChunkStreamID-64) % 256)
		buf[2] = byte(int(bh.ChunkStreamID-64) / 256)
		_, err := writer.Write(buf) // TODO: should check length?
		return err

	default:
		return errors.New("chunk stream id is out of range: %d must be in range [2, 65599]")
	}
}
