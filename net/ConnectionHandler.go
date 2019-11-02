package net

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"rsps/net/packet"
	"time"
)

var PACKET_SIZE = []int{0, 0, 0, 1, -1, 0, 0, 0, 0, 0, // 0
	0, 0, 0, 0, 8, 0, 6, 2, 2, 0, // 10
	0, 2, 0, 6, 0, 12, 0, 0, 0, 0, // 20
	0, 0, 0, 0, 0, 8, 4, 0, 0, 2, // 30
	2, 6, 0, 6, 0, -1, 0, 0, 0, 0, // 40
	0, 0, 0, 12, 0, 0, 0, 8, 8, 12, // 50
	8, 8, 0, 0, 0, 0, 0, 0, 0, 0, // 60
	6, 0, 2, 2, 8, 6, 0, -1, 0, 6, // 70
	0, 0, 0, 0, 0, 1, 4, 6, 0, 0, // 80
	0, 0, 0, 0, 0, 3, 0, 0, -1, 0, // 90
	0, 13, 0, -1, 0, 0, 0, 0, 0, 0, // 100
	0, 0, 0, 0, 0, 0, 0, 6, 0, 0, // 110
	1, 0, 6, 0, 0, 0, -1, 0, 2, 6, // 120
	0, 4, 6, 8, 0, 6, 0, 0, 0, 2, // 130
	0, 0, 0, 0, 0, 6, 0, 0, 0, 0, // 140
	0, 0, 1, 2, 0, 2, 6, 0, 0, 0, // 150
	0, 0, 0, 0, -1, -1, 0, 0, 0, 0, // 160
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 170
	0, 8, 0, 3, 0, 2, 0, 0, 8, 1, // 180
	0, 0, 12, 0, 0, 0, 0, 0, 0, 0, // 190
	2, 0, 0, 0, 0, 0, 0, 0, 4, 0, // 200
	4, 0, 0, 0, 7, 8, 0, 0, 10, 0, // 210
	0, 0, 0, 0, 0, 0, -1, 0, 6, 0, // 220
	1, 0, 0, 0, 6, 0, 6, 8, 1, 0, // 230
	0, 4, 0, 0, 0, 0, -1, 0, -1, 4, // 240
	0, 0, 6, 6, 0, 0, 0, // 250
}

type ConnectionHandler struct {
	//Connections map[string]*Connection
	LoginHandler *LoginHandler
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		LoginHandler: &LoginHandler{},
		//Connections: make(map[string]*Connection),
	}
}

const PORT = "43594"

func (c *ConnectionHandler) Listen() {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Listening on %s", PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		newConn := &Connection{
			Conn:    conn,
			Context: context.Background(),
		}
		newConn.SetValue(LoginState, 0)
		go c.handleConnection(newConn)
	}
}

func (c *ConnectionHandler) handleConnection(conn *Connection) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		time.Sleep(1 * time.Second)
		if conn.GetValue(LoginState) != 2 {
			c.LoginHandler.HandlePacket(conn)
			continue
		}

		buf := bufio.NewReader(conn)

		first, _ := buf.ReadByte()
		opCode := int(first & 0xff) // - int(conn.Decryptor.Rand()&0xff)

		var size int
		//log.Printf("opcode: %d", opCode)
		if opCode >= 0 && opCode < len(PACKET_SIZE) {
			size = PACKET_SIZE[opCode]
			log.Printf("opcode: %+v, size: %+v", opCode, size)

			//data := make([]byte, size)
			//
			//_, err := buf.Read(data)
			//if err != nil {
			//	log.Fatal(err.Error())
			//}
			//
			//log.Printf("data: %+v", data)
		}

		if conn.GetValue("INITIALIZED") != 1 {
			w := NewInterfaceText("100%", 21)

			conn.Write(w.Bytes())

			conn.W([]int{-48, -1, -1})
			conn.W([]int{73, 1, 18, 1, -110, 81, 0, 59, -26, -43, -80, 7, -1, 16, -52, 0, -1, -1, 0, 0, 0, 0, 1, 18, 0, 1, 26, 1, 36, 1, 0, 1, 33, 1, 42, 1, 10, 0, 0, 0, 0, 0, 3, 40, 3, 55, 3, 51, 3, 52, 3, 53, 3, 54, 3, 56, 0, 0, 1, -88, -5, 9, 73, 127, 3, 0, 0})

			var x uint16
			x = (1 << 8) + 18
			y := (1 << 8) + -110 + 256
			w = new(bytes.Buffer)
			_ = binary.Write(w, binary.BigEndian, &MapRegionPacket{
				Id: 73,
				X:  274,
				Y:  402,
			})
			log.Printf("map region: %d %d", x, y)
			conn.SetValue("INITIALIZED", 1)
		} else {
			playerUpdatePacket := packet.NewPlayerUpdatePacket().
				SetUpdateRequired(true).
				SetType(packet.Idle).
				Build()
			conn.Wb(playerUpdatePacket)
		}
	}
}

func NewInterfaceText(text string, interfaceId uint16) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.Write([]byte{126})
	buf.Write([]byte{0, byte(len(text) + 3)})
	buf.Write([]byte(text))
	buf.Write([]byte{10})
	buf.Write([]byte{byte(interfaceId << 8), byte(interfaceId)})
	log.Printf("Writing interface: %+v", buf.Bytes())
	return buf
}

type MapRegionPacket struct {
	Id byte
	X  uint16
	Y  uint16
}

type LocationPacket struct {
	Id   byte
	Mask byte
	Y    uint8
	X    uint8
}

type Packet struct {
	Opcode byte
	Size   byte
}

func (c *Connection) WriteFrame(id int, b []int) {
	buf := make([]byte, len(b))
	for k, v := range b {
		if v < 0 {
			buf[k] = byte(256 + v)
		} else {
			buf[k] = byte(v)
		}
	}

	log.Printf("Writing Packet id: %d, %+v, %s", id, buf, string(buf))
	buf = append([]byte{byte(id)}, buf...)
	c.Write(buf)
}

func (c *Connection) Wb(b []byte) {
	log.Printf("Writing Packet: %+v", b)
	c.Write(b)
}

func (c *Connection) W(b []int) {
	buf := make([]byte, len(b))
	for k, v := range b {
		if v < 0 {
			buf[k] = byte(256 + v)
		} else {
			buf[k] = byte(v)
		}
	}

	log.Printf("Writing Packet: %+v", buf)
	c.Write(buf)
}
