package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/gtank/isaac"
	"log"
	"net"
	"rsps/entity"
	"rsps/net/packet"
	"rsps/net/packet/incoming"
	"rsps/net/packet/outgoing"
	"time"
)

const (
	Disconnected   = -1
	HandshakeStage = 0
	LoginStage     = 1
	Initialize     = 2
	IngameStage    = 3
)

type TCPClient struct {
	World        *entity.World
	Player       *entity.Player
	loginState   int
	connection   net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	Upstream     chan entity.UpstreamMessage
	Downstream   chan entity.DownstreamMessage
	loginHandler *UpstreamLoginHandler
	Encryptor    *isaac.ISAAC
	Decryptor    *isaac.ISAAC
}

func NewTcpClient(connnection net.Conn, loginHandler *UpstreamLoginHandler, world *entity.World) *TCPClient {
	player := entity.NewPlayer()
	world.AddPlayerToRegion(player)
	return &TCPClient{
		World:        world,
		Player:       player,
		connection:   connnection,
		reader:       bufio.NewReader(connnection),
		writer:       bufio.NewWriter(connnection),
		Upstream:     make(chan entity.UpstreamMessage, 64),
		Downstream:   make(chan entity.DownstreamMessage, 256),
		loginHandler: loginHandler,
	}
}

type initType struct{}

func (i *initType) Write(writer *bufio.Writer) {
	binary.Write(writer, binary.BigEndian, []int{208, 255, 255})
}

func (client *TCPClient) Read() {
	defer client.connectionTerminated()

connectionLoop:
	for {
		switch client.loginState {
		case Disconnected:
			break connectionLoop
		case HandshakeStage, LoginStage:
			client.loginHandler.HandlePacket(client)
			break
		case Initialize:
			client.Enqueue(&initType{})
			client.Enqueue(&outgoing.MapRegionPacket{Position: client.Player.Position})
			client.Enqueue(outgoing.NewPlayerUpdatePacket(client.Player).SetUpdateRequired(true).SetTyp(outgoing.Teleport))
			client.Enqueue(&flush{})
			client.loginState = IngameStage
			break
		case IngameStage:
			opcode, err := client.reader.ReadByte()
			if err != nil {
				break connectionLoop
			}

			var size uint8
			if int(opcode) < len(PACKET_SIZE) {
				if PACKET_SIZE[opcode] == -1 {
					size, err = client.reader.ReadByte()
					if err != nil {
						log.Printf("error reading packetId %s", err.Error())
						continue
					}
				} else {
					size = uint8(PACKET_SIZE[opcode])
				}
			}

			p := &packet.Packet{
				Opcode: opcode,
				Size:   uint16(size),
			}
			if size > 0 {
				payload := make([]byte, size)
				_, err = client.reader.Read(payload)
				if err != nil {
					log.Printf("error reading payload %s", err.Error())
					continue
				}
				p.Buffer = bytes.NewBuffer(payload)
			}

			client.Upstream <- p
			break
		}
	}
}

func (client *TCPClient) Tick() {
	for {
		<-time.After(600 * time.Millisecond)
		if client.loginState == Disconnected {
			return
		}
		client.Player.Tick()
		for _, v := range client.Player.OutgoingQueue {
			client.Enqueue(v)
		}
		client.Player.OutgoingQueue = make([]entity.DownstreamMessage, 0)
		//
		//if client.Player.RegionChanged {
		//	client.Enqueue(&outgoing.MapRegionPacket{Position: client.Player.Position})
		//}
		client.Enqueue(outgoing.NewPlayerUpdatePacket(client.Player).SetUpdateRequired(true))
		client.Enqueue(&flush{})
	}
}

func (client *TCPClient) ProcessUpstream() {
	for upstreamMessage := range client.Upstream {
		if msg, ok := upstreamMessage.(*packet.Packet); ok {
			if !isIgnored(msg.Opcode) {
				//log.Printf("upstreamMessage: %+v", msg)
			}
			if msg.Opcode == 164 ||  msg.Opcode == 248{
				incoming.Packets[msg.Opcode].HandlePacket(client.Player, msg)
			}
		}
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

func (client *TCPClient) Write() {
	for downstreamMessage := range client.Downstream {
		switch downstreamMessage.(type) {
		default:
			downstreamMessage.Write(client.writer)
		case *flush:
			downstreamMessage.Write(client.writer)
			client.Player.PostUpdate()
		}
	}
}

func (client *TCPClient) connectionTerminated() {
	close(client.Downstream)
	close(client.Upstream)
	client.loginState = Disconnected
}

func (client *TCPClient) Enqueue(msg entity.DownstreamMessage) {
	client.Downstream <- msg
}

type flush struct{}

func (f *flush) Write(writer *bufio.Writer) {
	writer.Flush()
}

func (client *TCPClient) Flush() {
	client.Downstream <- &flush{}
}
