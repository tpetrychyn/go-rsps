package packet

import "rsps/entity"

type Packet struct {
	Opcode  byte
	Size    uint8
	Payload []byte
}

type PacketListener interface {
	HandlePacket(*entity.Player, *Packet)
}