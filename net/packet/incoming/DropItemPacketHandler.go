package incoming

import (
	"rsps/entity"
	"rsps/net/packet"
)

type DropItemPacketHandler struct {}

func (d *DropItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	itemId := packet.ReadShortA()
	_ = packet.ReadShort()
	slot := packet.ReadShortA()

	player.Inventory.DropItem(int(itemId), int(slot))
}
