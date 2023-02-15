package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"rtmp-parser/Handshake"
)

var RTMPVersion = 3

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", ":1935")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			return err
		}

		go func() {
			err := s.createConnector(conn)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (s *Server) createConnector(conn net.Conn) error {
	err := s.handshakeWithClient(conn)
	if err != nil {
		return err
	}

	c := NewConnector(conn)

	go c.RunRead()
	go c.RunWrite()

	return nil
}

func (s *Server) handshakeWithClient(conn net.Conn) error {
	// Recv C0
	c0, err := Handshake.DecodeS0C0(conn)
	if err != nil {
		return err
	}

	if c0.Version > byte(RTMPVersion) {
		return errors.New("version non supported")
	}

	s0 := Handshake.NewS0C0(byte(RTMPVersion))

	if err := s0.EncodeS0C0(conn); err != nil {
		return err
	}

	// Send S1
	s1 := Handshake.NewS1C1()

	if err := s1.EncodeS1C1(conn); err != nil {
		return err
	}

	// Recv C1
	c1, err := Handshake.DecodeS1C1(conn)
	if err != nil {
		return err
	}

	// Send S2
	s2 := Handshake.NewS2C2(c1.Time, c1.Random)
	if err := s2.EncodeS2C2(conn); err != nil {
		return err
	}

	c2, err := Handshake.DecodeS2C2(conn)
	if err != nil {
		return err
	}

	// Check random echo
	if !bytes.Equal(c2.Random[:], s1.Random[:]) {
		return errors.New("random echo is not matched")
	}

	return nil
}
