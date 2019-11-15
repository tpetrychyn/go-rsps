package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
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
	invItem := e.player.Inventory.FindItem(int(itemId))
	if invItem == nil {
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.InventoryItemPacket{Slot: int(invSlot), Item: model.NilItem})
		return
	}

	def := util.GetItemDefinition(int(itemId))
	slot := outgoing.EQUIPMENT_SLOTS[def.Equipment.Slot]
	e.player.UpdateFlag.SetAppearance()

	equippedItem := e.Items[slot]
	e.player.Inventory.Items[invSlot] = equippedItem
	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: int(invSlot),
		Item: equippedItem,
	})

	e.Items[slot] = invItem
	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: e.ItemContainer,
		InterfaceId:   model.EQUIPMENT_INTERFACE_ID,
	})
}

func (e *Equipment) RemoveItem(slot, id uint16) {
	if e.Items[slot].ItemId != int(id) {
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		e.Items[slot] = model.NilItem
		e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
			ItemContainer: e.ItemContainer,
			InterfaceId: model.EQUIPMENT_INTERFACE_ID,
		})
		return
	}

	e.Items[slot] = model.NilItem
	e.player.Inventory.AddItem(int(id), 1)
	e.SetItem(0, 0, int(slot))
	e.player.UpdateFlag.SetAppearance()

	e.player.OutgoingQueue = append(e.player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: e.ItemContainer,
		InterfaceId: model.EQUIPMENT_INTERFACE_ID,
	})
}


