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

func NewTcpClient(connection net.Conn, loginHandler *LoginHandler, world *entity.World) *TCPClient {
	player := world.AddPlayer()
	return &TCPClient{
		World:        world,
		Player:       player,
		connection:   connection,
		reader:       bufio.NewReader(connection),
		writer:       bufio.NewWriter(connection),
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
			client.Enqueue(&outgoing.PlayerDetailsPacket{Pid:client.Player.Id})
			client.Enqueue(&outgoing.SendSongPacket{Song:-1})
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

func (client *TCPClient) UpdatePacket(wg *sync.WaitGroup) {
	defer wg.Done()
	client.Enqueue(outgoing.NewPlayerUpdatePacket(client.Player))
	client.Enqueue(outgoing.NewNpcUpdatePacket(client.Player))
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
		if client == nil || client.writer == nil {
			log.Printf("write nil finder client %+v client.writer %+v", client, client.writer)
		}
		if downstreamMessage == nil {
			log.Printf("write nil finder downstreamMessage %+v", downstreamMessage)
		}
		downstreamMessage.Write(client.writer)
		switch downstreamMessage.(type) {
		case *flush:
			_ = client.writer.Flush()
			client.Player.PostUpdate()
		case *outgoing.LogoutPacket:
			_ = client.writer.Flush()
			return
		}
	}
}

func (client *TCPClient) connectionTerminated() {
	log.Printf("connection dropped %+v", client.Player)
	client.World.RemovePlayer(client.Player.Id)
	client.loginState = Disconnected
	close(client.Downstream)
	close(client.Upstream)
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
