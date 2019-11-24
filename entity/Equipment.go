package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/repository"
	"rsps/util"
)

type Equipment struct {
	player *Player
	*model.ItemContainer
}

func NewEquipment(player *Player) *Equipment {
	return &Equipment{
		player: player,
		ItemContainer: model.NewItemContainer(14),
	}
}

func (e *Equipment) EquipItem(invSlot, itemId uint16) {
	_, invItem := e.player.Inventory.FindItem(int(itemId))
	if invItem == nil {
		return
	}

	def := util.GetItemDefinition(int(itemId))
	slot := outgoing.EQUIPMENT_SLOTS[def.Equipment.Slot]
	e.player.UpdateFlag.SetAppearance()

	equippedItem := e.Items[slot]
	e.player.Inventory.SetItem(equippedItem.ItemId, equippedItem.Amount, int(invSlot))
	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.InterfaceItemPacket{
		InterfaceId: model.INVENTORY_INTERFACE_ID,
		Slot: int(invSlot),
		Item: equippedItem,
	})


	e.SetItem(invItem.ItemId, invItem.Amount, slot)
	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: e.ItemContainer,
		InterfaceId:   model.EQUIPMENT_INTERFACE_ID,
	})

	go repository.InventoryRepositorySingleton.Save(e.player.Name, e.player.Inventory.Items)
	go repository.EquipmentRepositorySingleton.Save(e.player.Name, e.Items)
}

func (e *Equipment) RemoveItem(slot, id uint16) {
	if e.Items[slot].ItemId != int(id) {
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		e.Items[slot] = &model.Item{}
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
			ItemContainer: e.ItemContainer,
			InterfaceId: model.EQUIPMENT_INTERFACE_ID,
		})
		return
	}

	e.Items[slot] = &model.Item{}
	e.player.Inventory.AddItem(int(id), 1)
	e.SetItem(0, 0, int(slot))
	e.player.UpdateFlag.SetAppearance()

	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: e.ItemContainer,
		InterfaceId: model.EQUIPMENT_INTERFACE_ID,
	})

	go repository.EquipmentRepositorySingleton.Save(e.player.Name, e.Items)
}


