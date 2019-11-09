package packet

import "rsps/model"

type Packet struct {
	Opcode  byte
	Size    uint8
	Payload []byte
}

type PacketListener interface {
	HandlePacket(*model.Player, *Packet)
}