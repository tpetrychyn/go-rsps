package incoming

import (
	"log"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

type RemoveItemPacketHandler struct {}

func (e *RemoveItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	interfaceId := packet.ReadShortA() //interfaceId?
	slot := packet.ReadShortA()
	id := packet.ReadShortA()

	switch interfaceId {
	case model.EQUIPMENT_INTERFACE_ID:
		player.Equipment.RemoveItem(slot, id)
	case model.BANK_INVENTORY_INTERFACE_ID:
		player.Bank.DepositItem(int(slot), int(id), 1)
	case model.BANK_INTERFACE_ID:
		player.Bank.WithdrawItem(int(slot), int(id), 1)
	default:
		log.Printf("remove item from interface %d", interfaceId)
	}

}
