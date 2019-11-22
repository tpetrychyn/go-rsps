package incoming

import (
	"log"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

type BankFivePacketHandler struct {}

func (b *BankFivePacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	interfaceId := packet.ReadLEShortA()
	itemId := packet.ReadLEShortA()
	slot := packet.ReadLEShort()

	switch interfaceId  {
	case model.BANK_INVENTORY_INTERFACE_ID:
		player.Bank.DepositItem(int(slot), int(itemId), 5)
	case model.BANK_INTERFACE_ID:
		player.Bank.WithdrawItem(int(slot), int(itemId), 5)
	default:
		log.Printf("bank five slot %d inter %d item %d", slot, interfaceId, itemId)
	}
}

type BankTenPacketHandler struct {}

func (b *BankTenPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	interfaceId := packet.ReadLEShort()
	itemId := packet.ReadShortA()
	slot := packet.ReadShortA()

	switch interfaceId  {
	case model.BANK_INVENTORY_INTERFACE_ID:
		player.Bank.DepositItem(int(slot), int(itemId), 10)
	case model.BANK_INTERFACE_ID:
		player.Bank.WithdrawItem(int(slot), int(itemId), 10)
	default:
		log.Printf("bank ten slot %d inter %d item %d", slot, interfaceId, itemId)
	}
}

type BankAllPacketHandler struct {}

func (b *BankAllPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	slot := packet.ReadShortA()
	interfaceId := packet.ReadShort()
	itemId := packet.ReadShortA()

	switch interfaceId  {
	case model.BANK_INVENTORY_INTERFACE_ID:
		player.Bank.DepositItem(int(slot), int(itemId), -1)
	case model.BANK_INTERFACE_ID:
		player.Bank.WithdrawItem(int(slot), int(itemId), -1)
	default:
		log.Printf("bank all slot %d inter %d item %d", slot, interfaceId, itemId)
	}
}
