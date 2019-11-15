package entity

import (
	"log"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Inventory struct {
	player *Player
	*model.ItemContainer
}

func NewInventory(player *Player) *Inventory {
	return &Inventory{
		player: player,
		ItemContainer: model.NewItemContainer(28),
	}
}

func (i *Inventory) AddItem(id, amount int) int {
	slot := -1
	for k, v := range i.Items {
		if v.ItemId == 0 {
			slot = k
			i.SetItem(id, amount, slot)
			break
		}

		if k == int(i.Capacity-1) {
			log.Printf("inventory is full")
		}
	}

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: slot,
		Item: &model.Item{
			ItemId: id,
			Amount: 1,
		},
	})
	return slot
}

func (i *Inventory) SwapItems(from, to int) {
	fromItem := i.Items[from]
	toItem := i.Items[to]

	if fromItem == nil {
		log.Printf("no item found in that slot")
		return
	}

	i.Items[to] = fromItem
	i.Items[from] = toItem

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: from,
		Item: toItem,
	})

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: to,
		Item: fromItem,
	})
}

func (i *Inventory) DropItem(itemId, slot int) {
	invItem := i.FindItem(itemId)
	if invItem == nil {
		i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		return
	}

	i.Items[slot] = model.NilItem
	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InventoryItemPacket{Slot: slot, Item: model.NilItem})
	i.player.Region.CreateGroundItemAtPosition(i.player, invItem, i.player.Position)
}
