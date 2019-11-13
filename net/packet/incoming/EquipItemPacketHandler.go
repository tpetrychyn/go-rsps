package incoming

import (
	"rsps/entity"
	"rsps/net/packet"
)

type EquipItemPacketHandler struct {}

func (e *EquipItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	itemId := packet.ReadShort()
	slot := packet.ReadShortA()
	_ = packet.ReadShortA() //interfaceId

	player.EquipItem(slot, itemId)
}
