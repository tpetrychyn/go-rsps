package incoming

import (
	"log"
	"rsps/entity"
	"rsps/net/packet"
)

type EquipItemPacketHandler struct {}

func (e *EquipItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	itemId := packet.ReadShort()
	slot := packet.ReadShortA()
	interfaceId := packet.ReadShortA()

	log.Printf("id %+v slot %+v interfaceId %+v", itemId, slot, interfaceId)
	player.EquipItem(slot, itemId)
}
