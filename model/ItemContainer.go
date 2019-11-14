package model

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
		items[k] = NilItem
	}
	return &ItemContainer{
		Capacity: capacity,
		Items: items,
	}
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