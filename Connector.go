package main

import (
	"fmt"
	"net"
	"rtmp-parser/Chunk"
)

type Connector struct {
	con   net.Conn
	write chan *Chunk.Chunk
	read  chan *Chunk.Chunk
}

func NewConnector(con net.Conn) *Connector {
	return &Connector{
		con:   con,
		write: make(chan *Chunk.Chunk),
		read:  make(chan *Chunk.Chunk),
	}
}

func (c *Connector) RunRead() {
	for {
		chunk, err := Chunk.ReadChunk(c.con)
		if err != nil {
			continue
		}

		chunk.Execute(c.write)
	}
}

func (c *Connector) RunWrite() {
	for chunk := range c.write {
		err := chunk.ChunkBasicHeader.EncodeChunkBasicHeader(c.con)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = chunk.ChunkMessageHeader.EncodeChunkMessageHeader(c.con)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = c.con.Write(chunk.Body[:])
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
