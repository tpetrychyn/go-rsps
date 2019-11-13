package model

import (
	"log"
)

const (
	INVENTORY_INTERFACE_ID = 3214
	EQUIPMENT_INTERFACE_ID = 1688
)

type ItemContainer struct {
	Capacity uint
	Items    []*Item
}

func NewItemContainer(capacity uint) *ItemContainer {
	items := make([]*Item, capacity)
	for k, _ := range items {
		items[k] = &Item{}
	}
	return &ItemContainer{
		Capacity: capacity,
		Items: items,
	}
}

func (i *ItemContainer) SwapItems(from, to int) {
	fromItem := i.Items[from]
	toItem := i.Items[to]

	if fromItem == nil {
		log.Printf("no item found in that slot")
		return
	}

	i.Items[to] = fromItem
	i.Items[from] = toItem
}

func (i *ItemContainer) AddItem(id, amount int) int {
	slot := -1
	for k, v := range i.Items {
		if v.ItemId == 0 {
			slot = k
			i.SetItem(id, amount, slot)
			break
		}

		if k == int(i.Capacity-1) {
			log.Printf("container is full")
		}
	}
	return slot
}

func (i *ItemContainer) SetItem(id, amount, slot int) {
	i.Items[slot] = &Item{
		ItemId: id,
		Amount: amount,
	}
}

func (i *ItemContainer) FindItem(id int) *Item {
	for _, v := range i.Items {
		if v.ItemId == id {
			return v
		}
	}
	return nil
}