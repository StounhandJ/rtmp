package Chunk

func createWindowAcknowledgementSize(size []byte, csId int) *Chunk {
	bh := NewChunkBasicHeader(0x1, csId)

	mh := NewChunkMessageHeader(0x1, 0x5, 0, 0, 4, 0)

	return NewChunk(bh, mh, size[:4])
}

func createSetPeerBandwidth(size []byte, limitType, csId int) *Chunk {
	buf := make([]byte, 5)

	copy(buf[:], size[:])
	buf[4] = byte(limitType)

	bh := NewChunkBasicHeader(0x0, csId)

	mh := NewChunkMessageHeader(0x0, 0x6, 0, 0, 5, 0)

	return NewChunk(bh, mh, buf[:])
}
