package model

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
