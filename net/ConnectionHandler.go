package net

import (
	"bufio"
	"fmt"
	"log"
	"rsps/net/packet"
	"rsps/net/packet/incoming"
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

func (c *ConnectionHandler) Writer(conn *Connection) {
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
			//outgoing.SendMapRegion(conn)
			conn.W([]int{81, 0, 59, -26, -43, -80, 7, -1, 16, -52, 0, -1, -1, 0, 0, 0, 0, 1, 18, 0, 1, 26, 1, 36, 1, 0, 1, 33, 1, 42, 1, 10, 0, 0, 0, 0, 0, 3, 40, 3, 55, 3, 51, 3, 52, 3, 53, 3, 54, 3, 56, 0, 0, 1, -88, -5, 9, 73, 127, 3, 0, 0})

			playerUpdatePacket := packet.NewPlayerUpdatePacket(conn.Player).
				SetUpdateRequired(true).
				Build()
			conn.Wb(playerUpdatePacket)

			conn.SetValue("INITIALIZED", 1)
		} else {
			conn.Player.Tick()

			playerUpdatePacket := packet.NewPlayerUpdatePacket(conn.Player).
				SetUpdateRequired(true).
				Build()
			conn.Wb(playerUpdatePacket)

			conn.Player.PostUpdate()

			for len(conn.PacketQueue) > 0 {
				p := conn.PacketQueue[0]
				conn.PacketQueue = conn.PacketQueue[1:]
				handler := incoming.Packets[p.Opcode]
				if handler != nil {
					handler.HandlePacket(conn.Player, p)
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (c *ConnectionHandler) Listener(conn *Connection) {
	fmt.Printf("Serving %s\n", conn.TCPConn.RemoteAddr().String())

	for {
		if conn.Player.LoginState != 2 {
			continue
		}

		buf := bufio.NewReader(conn.TCPConn)

		opCode, err := buf.ReadByte()
		if err != nil {
			log.Printf("error reading packetId %s", err.Error())
			conn.Player.LoginState = -1
			return
		}

		//opCode := byte(int(opCode & 0xff) - int(conn.Decryptor.Rand()&0xff))

		if isIgnored(opCode) { continue }

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
			Opcode:  opCode,
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

func isIgnored(opCode byte) bool {
	ignoredPackets := []byte{0, 3}
	for _, v := range ignoredPackets {
		if opCode == v {
			return true
		}
	}
	return false
}



