package entity

import (
	"errors"
	"log"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/repository"
	"rsps/util"
)

type Inventory struct {
	player *Player
	*model.ItemContainer
	CurrentInterfaceId int
}

func NewInventory(player *Player) *Inventory {
	return &Inventory{
		player:        player,
		ItemContainer: model.NewItemContainer(28),
		CurrentInterfaceId: model.INVENTORY_INTERFACE_ID,
	}
}

func (i *Inventory) AddItem(id, amount int) error {
	var slot = -1
	if !util.GetItemDefinition(id).Stackable && !util.GetItemDefinition(id).Noted && amount > 1 {
		var err error
		for a := 0; a < amount; a++ {
			err = i.AddItem(id, 1)
		}
		return err
	}

	// try to append stackable item
	if util.GetItemDefinition(id).Stackable || util.GetItemDefinition(id).Noted {
		for k, v := range i.Items {
			if v.ItemId == id {
				slot = k
				amount = v.Amount + amount
				i.SetItem(id, amount, slot)
				break
			}
		}
	}

	// if not existing/stackable, try filling a slot
	if slot == -1 {
		for k,v := range i.Items {
			// fill empty slot
			if v.ItemId == 0 {
				slot = k
				i.SetItem(id, amount, slot)
				break
			}
		}
	}

	// full
	if slot == -1 {
		i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Your inventory is too full to hold anymore."})
		return errors.New("inventory is full")
	}

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InterfaceItemPacket{
		InterfaceId: i.CurrentInterfaceId,
		Slot:        slot,
		Item: &model.Item{
			ItemId: id,
			Amount: amount,
		},
	})

	go repository.InventoryRepositorySingleton.Save(i.player.Name, i.Items)
	return nil
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

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InterfaceItemPacket{
		InterfaceId: i.CurrentInterfaceId,
		Slot:        from,
		Item:        toItem,
	})

	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InterfaceItemPacket{
		InterfaceId: i.CurrentInterfaceId,
		Slot:        to,
		Item:        fromItem,
	})

	go repository.InventoryRepositorySingleton.Save(i.player.Name, i.Items)
}

func (i *Inventory) DropItem(itemId, slot int) {
	_, invItem := i.FindItem(itemId)
	if invItem == nil {
		i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		return
	}

	i.Items[slot] = &model.Item{}
	i.player.OutgoingQueue = append(i.player.OutgoingQueue, &outgoing.InterfaceItemPacket{InterfaceId: model.INVENTORY_INTERFACE_ID, Slot: slot, Item: &model.Item{}})
	i.player.Region.CreateGroundItemAtPosition(i.player, invItem, i.player.Position)

	go repository.InventoryRepositorySingleton.Save(i.player.Name, i.Items)
}

func (i *Inventory) IsFull() bool {
	for _, v := range i.Items {
		if v.ItemId == 0 {
			return false
		}
	}
	return true
}
