package Chunk

import "io"

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
	buf := make([]byte, 1)

	//TODO сдвиг при большом csid

	buf[0] = bh.Fmt | byte(bh.ChunkStreamID)

	if _, err := writer.Write(buf[:]); err != nil {
		return err
	}
	return nil
}
