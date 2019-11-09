package outgoing

import (
	"rsps/model"
	"rsps/net/packet"
)

func SendMapRegion(p *model.Position) *packet.Packet {
	s := model.NewStream()
	s.WriteWord(uint(p.GetRegionX() + 6) + 128)
	s.WriteWord(uint(p.GetRegionY() + 6))

	return &packet.Packet{
		Opcode:  73,
		Size:    4,
		Payload: s.Flush(),
	}
}
