package entity

import (
	"github.com/google/uuid"
	"log"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Player struct {
	Id uuid.UUID
	*model.Movement
	MovementQueue   *MovementQueue
	Region          *Region
	Inventory       *model.ItemContainer
	Equipment       *model.ItemContainer
	OutgoingQueue   []DownstreamMessage
	LogoutRequested bool
}

var SIDEBARS = []int{2423, 3917, 638, 3213, 1644, 5608, 1151,
	18128, 5065, 5715, 2449, 904, 147, 962}

func NewPlayer() *Player {
	spawn := &model.Position{
		X: 3200,
		Y: 3200,
	}
	player := &Player{
		Id: uuid.New(), // TODO: Load this from database or something
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			SecondaryDirection: model.None,
			IsRunning:          true,
		},
	}
	player.MovementQueue = NewMovementQueue(player)

	for k, v := range SIDEBARS {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SidebarInterfacePacket{
			MenuId: k,
			Form:   v,
		})
	}

	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.ConfigurationPacket{
		InterfaceId: 173,
		State:       1,
	})

	inventory := model.NewItemContainer(28)
	inventory.Items[0] = &model.Item{
		ItemId: 995,
		Amount: 10000,
	}
	inventory.Items[1] = &model.Item{
		ItemId: 1351,
		Amount: 1,
	}
	inventory.Items[2] = &model.Item{
		ItemId: 579,
		Amount: 1,
	}
	player.Inventory = inventory
	player.Equipment = model.NewItemContainer(14)

	for k, v := range player.Inventory.Items {
		log.Printf("k %+v v %+v", k, v)
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.InventoryItemPacket{
			Slot: k,
			Item: v,
		})
	}

	return player
}

func (p *Player) PostUpdate() {
	p.PrimaryDirection = model.None
	p.SecondaryDirection = model.None
	p.LastDirection = model.None
}

func (p *Player) Tick() {
	p.MovementQueue.Tick()
}

func (p *Player) EquipItem(equipSlot, itemId uint16) {
	var invItem *model.Item
	var invSlot int
	for k, v := range p.Inventory.Items {
		if v.ItemId == int(itemId) {
			invItem = v
			invSlot = k
		}
	}
	if invItem == nil {
		log.Printf("you do not have that item")
		return
	}

	p.Equipment.Items[equipSlot] = p.Inventory.Items[invSlot]
	p.Inventory.Items[invSlot] = &model.Item{}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: invSlot,
		Item: &model.Item{},
	})
}

func (p *Player) GetPosition() *model.Position {
	return p.Position
}

func (p *Player) SetPosition(position *model.Position) {
	p.Position = position
}

func (p *Player) GetLastKnownRegion() *model.Position {
	return p.LastKnownRegion
}

func (p *Player) SetLastKnownRegion(position *model.Position) {
	p.LastKnownRegion = position
}

func (p *Player) GetPrimaryDirection() model.Direction {
	return p.PrimaryDirection
}

func (p *Player) SetPrimaryDirection(direction model.Direction) {
	p.PrimaryDirection = direction
}

func (p *Player) GetSecondaryDirection() model.Direction {
	return p.SecondaryDirection
}

func (p *Player) SetSecondaryDirection(direction model.Direction) {
	p.SecondaryDirection = direction
}

func (p *Player) GetLastDirection() model.Direction {
	return p.LastDirection
}

func (p *Player) SetLastDirection(direction model.Direction) {
	p.LastDirection = direction
}

func (p *Player) GetEquipment() *model.ItemContainer {
	return p.Equipment
}
