package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/repository"
)

const MAX_BANK_SIZE = 352 // Client crashes if larger

type Bank struct {
	player *Player
	*model.ItemContainer
}

func NewBank(player *Player) *Bank {
	return &Bank{
		player:        player,
		ItemContainer: model.NewItemContainer(MAX_BANK_SIZE),
	}
}

func (b *Bank) OpenBank() {
	b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: b.player.Inventory.ItemContainer,
		InterfaceId:   model.BANK_INVENTORY_INTERFACE_ID,
	})
	b.player.OutgoingQueue = append(b.player.OutgoingQueue, outgoing.NewInventoryInterfacePacket(model.BANK_WITH_INVENTORY_INTERFACE_ID, 5063))
	b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: b.ItemContainer,
		InterfaceId:   model.BANK_INTERFACE_ID,
	})
}

func (b *Bank) DepositItem(invSlot, itemId, requestedAmount int) {
	// TODO: multiples of non-stackable items
	// TODO: noted items
	invSlot, invItem := b.player.Inventory.FindItem(itemId)
	if invItem == nil {
		return
	}

	// bank all or bank max
	if requestedAmount == -1 || invItem.Amount < requestedAmount {
		requestedAmount = invItem.Amount
	}

	//def := util.GetItemDefinition(int(itemId))
	var slot = -1
	var amount int
	for k, v := range b.player.Bank.Items {
		// exists in bank
		if v.ItemId == itemId {
			amount = v.Amount + requestedAmount
			slot = k
			b.SetItem(itemId, amount, slot)
			break
		}
	}
	if slot == -1 {
		for k, v := range b.Items {
			if v.ItemId == 0 {
				slot = k
				amount = requestedAmount
				b.SetItem(itemId, requestedAmount, slot)
				break
			}
		}
	}

	if slot == -1 {
		b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Your bank is too full to hold anymore items."})
		return
	}

	b.player.Inventory.Items[invSlot].Amount -= requestedAmount
	if b.player.Inventory.Items[invSlot].Amount <= 0 {
		b.player.Inventory.Items[invSlot] = &model.Item{}
	}
	b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.InterfaceItemPacket{InterfaceId: model.BANK_INVENTORY_INTERFACE_ID, Slot: invSlot, Item: b.player.Inventory.Items[invSlot]})

	b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.InterfaceItemPacket{
		InterfaceId: model.BANK_INTERFACE_ID,
		Slot:        slot,
		Item: &model.Item{
			ItemId: itemId,
			Amount: amount,
		},
	})

	go repository.InventoryRepositorySingleton.Save(b.player.Name, b.player.Inventory.Items)
	go repository.BankRepositorySingleton.Save(b.player.Name, b.Items)
}

func (b *Bank) WithdrawItem(bankSlot, itemId, requestedAmount int) {
	_, bankItem := b.FindItem(itemId)
	if bankItem == nil {
		return
	}
	if requestedAmount == -1 || bankItem.Amount < requestedAmount {
		requestedAmount = bankItem.Amount
	}

	b.player.Inventory.CurrentInterfaceId = model.BANK_INVENTORY_INTERFACE_ID
	err := b.player.Inventory.AddItem(itemId, requestedAmount)
	if err != nil {
		return
	}

	b.Items[bankSlot].Amount -= requestedAmount
	if b.Items[bankSlot].Amount <= 0 {
		b.Items[bankSlot] = &model.Item{}
	}
	b.player.OutgoingQueue = append(b.player.OutgoingQueue, &outgoing.InterfaceItemPacket{InterfaceId: model.BANK_INTERFACE_ID, Slot: bankSlot, Item: b.Items[bankSlot]})

	go repository.BankRepositorySingleton.Save(b.player.Name, b.Items)
}
