package incoming

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

type MoveItemPacketHandler struct {}

func (e *MoveItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	interfaceId := packet.ReadLEShortA()
	_ = ^packet.ReadByte() == 1 //insertMode
	from := packet.ReadLEShortA()
	to := packet.ReadLEShort()

	switch interfaceId {
	case model.INVENTORY_INTERFACE_ID:
		player.Inventory.SwapItems(int(from), int(to))
	}
}
