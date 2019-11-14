package entity

import (
	"github.com/google/uuid"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Player struct {
	*model.Movement

	Id              uuid.UUID
	Name            string
	MovementQueue   *MovementQueue
	Region          *Region
	Inventory       *Inventory
	Equipment       *Equipment
	OutgoingQueue   []DownstreamMessage
	UpdateFlag      *model.UpdateFlag
	DelayedPacket   func()
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
		Id:         uuid.New(), // TODO: Load this from database or something
		UpdateFlag: &model.UpdateFlag{},
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			SecondaryDirection: model.None,
			IsRunning:          true,
		},
	}
	player.MovementQueue = NewMovementQueue(player)
	player.Inventory = NewInventory(player)
	player.Equipment = NewEquipment(player)

	return player
}

func (p *Player) LoadPlayer(name string) error {
	p.Name = name
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Inventory.ItemContainer,
		InterfaceId:   3214,
	})
	for k, v := range SIDEBARS {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SidebarInterfacePacket{
			MenuId: k,
			Form:   v,
		})
	}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Welcome to TaylorScape"})
	return nil
}

func (p *Player) Tick() {
	p.MovementQueue.Tick()
}

func (p *Player) PostUpdate() {
	p.PrimaryDirection = model.None
	p.SecondaryDirection = model.None
	p.LastDirection = model.None
}

func (p *Player) Teleport(position *model.Position) {
	p.LastDirection = model.None
	p.PrimaryDirection = model.None
	p.Position = position
	p.LastKnownRegion = p.Position
	p.MovementQueue.Clear()
	p.UpdateFlag.NeedsPlacement = true
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.MapRegionPacket{Position: p.Position})
}

func (p *Player) GetEquipmentItemContainer() *model.ItemContainer {
	return p.Equipment.ItemContainer
}

func (p *Player) GetUpdateFlag() *model.UpdateFlag {
	return p.UpdateFlag
}
