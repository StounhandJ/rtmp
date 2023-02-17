package Chunk

import (
	"encoding/binary"
	"fmt"
	"github.com/StounhandJ/go-amf"
	"io"
)

type Chunk struct {
	ChunkBasicHeader   *BasicHeader
	ChunkMessageHeader *MessageHeader
	Body               []byte
}

func NewChunk(bh *BasicHeader, mh *MessageHeader, body []byte) *Chunk {
	return &Chunk{
		ChunkBasicHeader:   bh,
		ChunkMessageHeader: mh,
		Body:               body,
	}
}

func ReadChunk(r io.Reader) (*Chunk, error) {
	bh, err := DecodeChunkBasicHeader(r)
	if err != nil {
		return nil, err
	}

	mh, err := decodeChunkMessageHeader(r, bh.Fmt)
	if err != nil {
		return nil, err
	}

	body := make([]byte, mh.MessageLength)

	if _, err := io.ReadAtLeast(r, body, int(mh.MessageLength)); err != nil {
		return nil, err
	}

	return &Chunk{
		ChunkBasicHeader:   bh,
		ChunkMessageHeader: mh,
		Body:               body[:],
	}, nil
}

func (c *Chunk) Execute(write chan<- *Chunk) {
	switch c.ChunkMessageHeader.MessageTypeID {
	case 1:
		// SetChunkSize
		var value uint32
		value |= uint32(c.Body[3])
		value |= uint32(c.Body[2]) << 8
		value |= uint32(c.Body[1]) << 16
		value |= uint32(c.Body[0]&0x7f) << 24
		break
	case 2:
		// AbortMessage
		var chunkStreamId uint32
		chunkStreamId = binary.BigEndian.Uint32(c.Body)
		fmt.Println(chunkStreamId)
		break
	case 17, 20:
		// AudioMessage, VideoMessage

		f1, s1, _ := amf.DecodeAMF0(c.Body)

		_, s2, _ := amf.DecodeAMF0(c.Body[s1:])

		_, _, _ = amf.DecodeAMF0(c.Body[s1+s2:])

		if f1 == "connect" {
			size := make([]byte, 4)
			binary.BigEndian.PutUint32(size[:], 4096)
			//TODO проверить правильность данных через свой декодер
			chWAL := createWindowAcknowledgementSize(size[:], c.ChunkBasicHeader.ChunkStreamID)

			write <- chWAL

			//chSPB := createSetPeerBandwidth(size[:], 1, c.ChunkBasicHeader.ChunkStreamID)
			//
			//write <- chSPB

			//bs := make([]byte, 4)
			//binary.LittleEndian.PutUint32(bs, 4048)
			//h1, err := w.Write(bs[:])
			//if err != nil {
			//	fmt.Println(err)
			//}
			//
			//bs2 := make([]byte, 5)
			//binary.LittleEndian.PutUint32(bs2, 4048)
			//bs2[4] = 0x01
			////h2, err := w.Write(bs2[:])
			//if err != nil {
			//	fmt.Println(err)
			//}
			////io.CopyN()
			//fmt.Println(h1, h2)
		}
		break
	}

}
