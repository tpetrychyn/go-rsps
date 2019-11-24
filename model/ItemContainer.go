package model

const (
	BANK_INVENTORY_INTERFACE_ID      = 5064
	BANK_WITH_INVENTORY_INTERFACE_ID = 5292
	BANK_INTERFACE_ID                = 5382
	INVENTORY_INTERFACE_ID           = 3214
	EQUIPMENT_INTERFACE_ID           = 1688
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
		Items:    items,
	}
}

func (i *ItemContainer) SetItem(id, amount, slot int) {
	i.Items[slot] = &Item{
		ItemId: id,
		Amount: amount,
	}
}

func (i *ItemContainer) FindItem(id int) (int, *Item) {
	for k, v := range i.Items {
		if v.ItemId == id {
			return k, v
		}
	}
	return -1, nil
}
