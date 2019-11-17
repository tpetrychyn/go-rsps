package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/gtank/isaac"
	"log"
	"net"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/incoming"
	"rsps/net/packet/outgoing"
	"sync"
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
	loginHandler *LoginHandler
	Encryptor    *isaac.ISAAC
	Decryptor    *isaac.ISAAC
}

func NewTcpClient(connnection net.Conn, loginHandler *LoginHandler, world *entity.World) *TCPClient {
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
			client.Player.UpdateFlag.SetAppearance()
			client.Player.UpdateFlag.NeedsPlacement = true
			client.Enqueue(outgoing.NewPlayerUpdatePacket(client.Player))
			for _, v := range client.Player.OutgoingQueue {
				client.Enqueue(v)
			}
			client.Player.OutgoingQueue = make([]entity.DownstreamMessage, 0)
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
				p.Payload = payload
				p.Buffer = bytes.NewBuffer(payload)
			}

			client.Upstream <- p
			break
		}
	}
}

func (client *TCPClient) Tick(wg *sync.WaitGroup) {
	defer wg.Done()
	client.Player.Tick()
	for _, v := range client.Player.OutgoingQueue {
		client.Enqueue(v)
	}
	client.Player.OutgoingQueue = make([]entity.DownstreamMessage, 0)

	if client.Player.LogoutRequested {
		client.loginState = Disconnected
	}
}

func (client *TCPClient) UpdatePacket() {
	client.Enqueue(outgoing.NewPlayerUpdatePacket(client.Player))
	client.Enqueue(&flush{})
}

func (client *TCPClient) ProcessUpstream() {
	for upstreamMessage := range client.Upstream {
		if msg, ok := upstreamMessage.(*packet.Packet); ok {
			if !isIgnored(msg.Opcode) && incoming.Packets[msg.Opcode] == nil {
				log.Printf("upstreamMessage: opcode %+v, size %+v, payload %+v \n", msg.Opcode, msg.Size, msg.Payload)
			}
			handler := incoming.Packets[msg.Opcode]
			if handler != nil {
				incoming.Packets[msg.Opcode].HandlePacket(client.Player, msg)
			}
		}
	}
}

func isIgnored(opCode byte) bool {
	// 0 keepalive
	// 241 click
	ignoredPackets := []byte{0, 3, 36, 77, 78, 86, 121, 210, 136, 241}
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
		case *outgoing.LogoutPacket:
			client.Player.Position = &model.Position{X: 0, Y: 0, Z: 255} // cheap hack to get other clients to delete this player..
			downstreamMessage.Write(client.writer)
			_ = client.writer.Flush()
			return
		case *flush:
			downstreamMessage.Write(client.writer)
			client.Player.PostUpdate()
		default:
			downstreamMessage.Write(client.writer)
		}
	}
}

func (client *TCPClient) connectionTerminated() {
	log.Printf("connection dropped %+v", client.Player)
	close(client.Downstream)
	close(client.Upstream)
	client.loginState = Disconnected
	client.Player.Position = &model.Position{X: 0, Y: 0, Z: 255}
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
