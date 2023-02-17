package main

import (
	"bufio"
	"fmt"
	"net"
	"rtmp-parser/Chunk"
)

type Connector struct {
	con    net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
	write  chan *Chunk.Chunk
	read   chan *Chunk.Chunk
}

func NewConnector(con net.Conn) *Connector {
	return &Connector{
		con:    con,
		writer: bufio.NewWriter(con),
		reader: bufio.NewReader(con),
		write:  make(chan *Chunk.Chunk),
		read:   make(chan *Chunk.Chunk),
	}
}

func (c *Connector) RunRead() {
	for {
		chunk, err := Chunk.ReadChunk(c.reader)
		if err != nil {
			continue
		}

		chunk.Execute(c.write)
	}
}

func (c *Connector) RunWrite() {
	for chunk := range c.write {
		err := chunk.ChunkBasicHeader.EncodeChunkBasicHeader(c.writer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = chunk.ChunkMessageHeader.EncodeChunkMessageHeader(c.writer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = c.writer.Write(chunk.Body[:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = c.writer.Flush()
		if err != nil {
			return
		}
	}
}
