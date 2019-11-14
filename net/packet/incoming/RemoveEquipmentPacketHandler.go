package incoming

import (
	"rsps/entity"
	"rsps/net/packet"
)

type RemoveEquipmentPacketHandler struct {}

func (e *RemoveEquipmentPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	_ = packet.ReadLEShortA() //interfaceId?
	slot := packet.ReadShortA()
	id := packet.ReadShortA()

	player.Equipment.RemoveItem(slot, id)
}
