package Chunk

import (
	"encoding/binary"
	"errors"
	"io"
	"rtmp-parser/utils"
)

type MessageHeader struct {
	Fmt             byte
	Timestamp       uint32 // fmt = 0
	TimestampDelta  uint32 // fmt = 1 | 2
	MessageLength   uint32 // fmt = 0 | 1
	MessageTypeID   byte   // fmt = 0 | 1
	MessageStreamID uint32 // fmt = 0
}

func NewChunkMessageHeader(fmt, messageTypeID byte, timestamp, timestampDelta, messageLength, messageStreamID uint32) *MessageHeader {
	return &MessageHeader{
		Fmt:             fmt,
		Timestamp:       timestamp,
		TimestampDelta:  timestampDelta,
		MessageLength:   messageLength,
		MessageTypeID:   messageTypeID,
		MessageStreamID: messageStreamID,
	}
}

func decodeChunkMessageHeader(r io.Reader, fmt byte) (*MessageHeader, error) {
	mh := &MessageHeader{
		Fmt: fmt,
	}

	buf := make([]byte, 11)
	cache32bits := make([]byte, 4)

	switch fmt {
	case 0:
		if _, err := io.ReadAtLeast(r, buf[:11], 11); err != nil {
			return nil, err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.Timestamp = binary.BigEndian.Uint32(cache32bits)
		copy(cache32bits[1:], buf[3:6]) // 24bits BE
		mh.MessageLength = binary.BigEndian.Uint32(cache32bits)
		mh.MessageTypeID = buf[6]                                  // 8bits
		mh.MessageStreamID = binary.LittleEndian.Uint32(buf[7:11]) // 32bits
		if mh.Timestamp == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return nil, err
			}
			mh.Timestamp = binary.LittleEndian.Uint32(cache32bits)
		}

	case 1:
		if _, err := io.ReadAtLeast(r, buf[:7], 7); err != nil {
			return nil, err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.TimestampDelta = binary.BigEndian.Uint32(cache32bits)
		copy(cache32bits[1:], buf[3:6]) // 24bits BE
		mh.MessageLength = binary.BigEndian.Uint32(cache32bits)
		mh.MessageTypeID = buf[6] // 8bits
		if mh.TimestampDelta == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return nil, err
			}
			mh.TimestampDelta = binary.LittleEndian.Uint32(cache32bits)
		}

	case 2:
		if _, err := io.ReadAtLeast(r, buf[:3], 3); err != nil {
			return nil, err
		}

		copy(cache32bits[1:], buf[0:3]) // 24bits BE
		mh.TimestampDelta = binary.BigEndian.Uint32(cache32bits)

		if mh.TimestampDelta == 0xffffff {
			_, err := io.ReadAtLeast(r, cache32bits, 4)
			if err != nil {
				return nil, err
			}
			mh.TimestampDelta = binary.LittleEndian.Uint32(cache32bits)
		}

	case 3:
		// DO NOTHING
	default:
		return nil, errors.New("unexpected fmt")
	}

	return mh, nil
}

func (mh *MessageHeader) EncodeChunkMessageHeader(writer io.Writer) error {
	mes := make([]byte, 15)

	switch mh.Fmt {
	case 0:
		copy(mes[3:6], utils.Uint32ToByteArray(mh.MessageLength, 3))
		mes[6] = mh.MessageTypeID
		copy(mes[7:11], utils.Uint32ToByteArray(mh.MessageStreamID, 4))
		if mh.Timestamp > 0xffffff {
			copy(mes[0:3], utils.Uint32ToByteArray(0xffffff, 3))
			copy(mes[10:14], utils.Uint32ToByteArray(mh.Timestamp, 4))
			_, err := writer.Write(mes[:14])
			if err != nil {
				return err
			}
		} else {
			copy(mes[0:3], utils.Uint32ToByteArray(mh.Timestamp, 3))
			_, err := writer.Write(mes[:11])
			if err != nil {
				return err
			}
		}
	case 1:
		copy(mes[3:6], utils.Uint32ToByteArray(mh.MessageLength, 3))
		mes[6] = mh.MessageTypeID
		if mh.TimestampDelta > 0xffffff {
			copy(mes[0:3], utils.Uint32ToByteArray(0xffffff, 3))
			copy(mes[7:11], utils.Uint32ToByteArray(mh.TimestampDelta, 4))
			_, err := writer.Write(mes[:10])
			if err != nil {
				return err
			}
		} else {
			copy(mes[0:3], utils.Uint32ToByteArray(mh.TimestampDelta, 3))
			_, err := writer.Write(mes[:7])
			if err != nil {
				return err
			}
		}

	case 2:
		binary.BigEndian.PutUint32(mes[0:3], mh.TimestampDelta)
		if mh.TimestampDelta > 0xffffff {
			binary.BigEndian.PutUint32(mes[0:3], 0xffffff)
			binary.BigEndian.PutUint32(mes[3:7], mh.TimestampDelta)
			_, err := writer.Write(mes[:7])
			if err != nil {
				return err
			}
		} else {
			binary.BigEndian.PutUint32(mes[0:3], mh.TimestampDelta)
			_, err := writer.Write(mes[:3])
			if err != nil {
				return err
			}
		}

	case 3:
		// DO NOTHING
	default:
		return errors.New("unexpected fmt")
	}
	return nil
}
