package net

import (
	"bufio"
	"bytes"
	"context"
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
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":43594")
	ln, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Listening on %s", PORT)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// TODO: Handle error
		}
		newConn := &Connection{
			TCPConn: conn,
			Context: context.Background(),
		}
		newConn.SetValue(LoginState, 0)
		go c.handleConnection(newConn)
		go c.PlayerTick(newConn)
	}
}

func (c *ConnectionHandler) PlayerTick(conn *Connection) {
	for {
		if conn.GetValue(LoginState) != 2 {
			c.LoginHandler.HandlePacket(conn)
			continue
		}

		if conn.GetValue("INITIALIZED") != 1 {
			w := NewInterfaceText("100%", 21)

			conn.Write(w.Bytes())

			conn.W([]int{-48, -1, -1})
			conn.W([]int{73, 1, 18, 1, -110, 81, 0, 59, -26, -43, -80, 7, -1, 16, -52, 0, -1, -1, 0, 0, 0, 0, 1, 18, 0, 1, 26, 1, 36, 1, 0, 1, 33, 1, 42, 1, 10, 0, 0, 0, 0, 0, 3, 40, 3, 55, 3, 51, 3, 52, 3, 53, 3, 54, 3, 56, 0, 0, 1, -88, -5, 9, 73, 127, 3, 0, 0})

			playerUpdatePacket := packet.NewPlayerUpdatePacket().
				SetUpdateRequired(true).
				SetType(packet.Idle).
				Build()
			conn.Wb(playerUpdatePacket)

			conn.SetValue("INITIALIZED", 1)
		} else {
			playerUpdatePacket := packet.NewPlayerUpdatePacket().
				SetUpdateRequired(true).
				SetType(packet.Idle).
				Build()
			conn.Wb(playerUpdatePacket)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (c *ConnectionHandler) handleConnection(conn *Connection) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		if conn.GetValue(LoginState) != 2 {
			continue
		}
		buf := bufio.NewReader(conn)

		id, err := buf.ReadByte()
		if err != nil {
			log.Printf("error reading packetId %s", err.Error())
		}

		opCode := int(id & 0xff) // - int(conn.Decryptor.Rand()&0xff)

		var size int
		if opCode >= 0 && opCode < len(PACKET_SIZE) {
			size = PACKET_SIZE[opCode]
		}

		if size > 0 {
			payload := make([]byte, size)
			_, err = buf.Read(payload)
			if err != nil {
				log.Printf("error reading payload %s", err.Error())
			}
			packet := []byte{id, byte(size)}
			packet = append(packet, payload...)
			conn.packetQueue.Write(packet)
			log.Printf("Read Packet id: %d, size: %d, payload: %+v", id, size, payload)
		} else {
			log.Printf("Read Packet id: %d, size: %d", id, size)
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
	//log.Printf("Writing Packet: %+v", b)
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
