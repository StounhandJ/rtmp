package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"net"
	"time"
)
import "fmt"

// требуется только ниже для обработки примера

func main() {

	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":1935")

	// Открываем порт
	conn, _ := ln.Accept()
	HandshakeWithClient2(conn, conn, &Config{})

	// Запускаем цикл
	for {
		f := bufio.NewReader(conn)

		var bh chunkBasicHeader
		if err := decodeChunkBasicHeader(f, make([]byte, 64*1024), &bh); err != nil {
			fmt.Println("Error")
		}

		var mh chunkMessageHeader
		if err := decodeChunkMessageHeader(f, bh.fmt, make([]byte, 64*1024), &mh); err != nil {
			fmt.Println("Error")
		}
		//// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadByte()
		////data := binary.BigEndian.Uint64(num)
		//fmt.Println(num)
		//_, err := conn.Write([]byte{3})
		//if err != nil {
		//	return
		//}
		//
		//message, _ = bufio.NewReader(conn).ReadByte()
		//// Распечатываем полученое сообщение
		fmt.Print("Message Received:", string(message))
		//// Процесс выборки для полученной строки
		////newmessage := strings.ToUpper(message)
		//// Отправить новую строку обратно клиенту
		//conn.Write([]byte(newmessage + "\n"))
	}
}

func HandshakeWithClient2(r io.Reader, w io.Writer, config *Config) error {
	d := NewDecoder(r)
	e := NewEncoder(w)

	// Recv C0
	var c0 S0C0
	if err := d.DecodeS0C0(&c0); err != nil {
		return err
	}

	// TODO: check c0 RTMP version

	// Send S0
	s0 := S0C0(RTMPVersion)
	if err := e.EncodeS0C0(&s0); err != nil {
		return err
	}

	// Send S1
	s1 := S1C1{
		Time: uint32(timeNow().UnixNano() / int64(time.Millisecond)),
	}
	copy(s1.Version[:], Version[:])
	if _, err := rand.Read(s1.Random[:]); err != nil { // Random Seq
		return err
	}
	if err := e.EncodeS1C1(&s1); err != nil {
		return err
	}

	// Recv C1
	var c1 S1C1
	if err := d.DecodeS1C1(&c1); err != nil {
		return err
	}

	// TODO: check c1 Client version. e.g. [9 0 124 2]

	// Send S2
	s2 := S2C2{
		Time:  c1.Time,
		Time2: uint32(timeNow().UnixNano() / int64(time.Millisecond)),
	}
	copy(s2.Random[:], c1.Random[:]) // echo c1 random
	if err := e.EncodeS2C2(&s2); err != nil {
		return err
	}

	// Recv C2
	var c2 S2C2
	if err := d.DecodeS2C2(&c2); err != nil {
		return err
	}

	if config.SkipHandshakeVerification {
		return nil
	}

	// Check random echo
	if !bytes.Equal(c2.Random[:], s1.Random[:]) {
		return errors.New("Random echo is not matched")
	}

	return nil
}

//package main
//
//import (
//	"io"
//	"net"
//
//	log "github.com/sirupsen/logrus"
//	"github.com/yutopp/go-rtmp"
//)
//
//func main() {
//	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1935")
//	if err != nil {
//		log.Panicf("Failed: %+v", err)
//	}
//
//	listener, err := net.ListenTCP("tcp", tcpAddr)
//	if err != nil {
//		log.Panicf("Failed: %+v", err)
//	}
//
//	srv := rtmp.NewServer(&rtmp.ServerConfig{
//		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
//			l := log.StandardLogger()
//			//l.SetLevel(logrus.DebugLevel)
//
//			return conn, &rtmp.ConnConfig{
//
//				ControlState: rtmp.StreamControlStateConfig{
//					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
//				},
//
//				Logger: l,
//			}
//		},
//	})
//	if err := srv.Serve(listener); err != nil {
//		log.Panicf("Failed: %+v", err)
//	}
//}
