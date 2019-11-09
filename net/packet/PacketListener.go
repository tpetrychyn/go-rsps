package packet

import (
	"rsps/entity"
)

type PacketListener interface {
	HandlePacket(*entity.Player, *Packet)
}
