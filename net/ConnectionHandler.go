package net

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"rsps/entity"
	"rsps/net/packet"
	"rsps/net/packet/handler"
	"time"
)

var PACKET_SIZE = []int8{0, 0, 0, 1, -1, 0, 0, 0, 0, 0, // 0
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
	LoginHandler *LoginHandler
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		LoginHandler: &LoginHandler{},
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
			Player: entity.NewPlayer(),
		}
		newConn.SetValue(LoginState, 0)
		go c.listener(newConn)
		go c.writer(newConn)
	}
}

func (c *ConnectionHandler) writer(conn *Connection) {
	for {
		if conn.GetValue(ConnectionStatus) == Disconnected {
			return
		}

		if conn.GetValue(LoginState) != 2 {
			c.LoginHandler.HandlePacket(conn)
			continue
		}

		if conn.GetValue("INITIALIZED") != 1 {

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


			for len(conn.PacketQueue) > 0 {
				p := conn.PacketQueue[0]
				conn.PacketQueue = conn.PacketQueue[1:]
				handler := handler.IncomingPackets[p.Opcode]
				if handler != nil {
					handler.HandlePacket(conn.Player, p)
				}
			}

			conn.Player.Tick()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (c *ConnectionHandler) listener(conn *Connection) {
	fmt.Printf("Serving %s\n", conn.TCPConn.RemoteAddr().String())

	for {
		if conn.GetValue(LoginState) != 2 {
			continue
		}

		buf := bufio.NewReader(conn.TCPConn)

		op, err := buf.ReadByte()
		if err != nil {
			log.Printf("error reading packetId %s", err.Error())
			conn.SetValue(ConnectionStatus, Disconnected)
			return
		}

		//opCode := int(op & 0xff) - int(conn.Decryptor.Rand()&0xff)
		opCode := op

		var ignored bool
		for _, v := range IGNORED_PACKETS {
			if opCode == v {
				ignored = true
				break
			}
		}
		if ignored { continue }

		var size uint8
		if int(opCode) < len(PACKET_SIZE) {
			if PACKET_SIZE[opCode] == -1 {
				size, err = buf.ReadByte()
				if err != nil {
					log.Printf("error reading packetId %s", err.Error())
					continue
				}
			} else {
				size = uint8(PACKET_SIZE[opCode])
			}
		}

		p := &packet.Packet{
			Opcode:  op,
			Size:    size,
		}
		if size > 0 {
			payload := make([]byte, size)
			_, err = buf.Read(payload)
			if err != nil {
				log.Printf("error reading payload %s", err.Error())
				continue
			}
			p.Payload = payload
		}
		log.Printf("Read Packet id: %d, size: %d, payload: %+v", p.Opcode, p.Size, p.Payload)
		conn.PacketQueue = append(conn.PacketQueue, p)
	}
}

var IGNORED_PACKETS = []byte{0, 3}


