package entity

import (
	"github.com/google/uuid"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Player struct {
	Id                 uuid.UUID
	MovementQueue      *MovementQueue
	Position           *model.Position
	Region             *Region
	LastKnownRegion    *model.Position
	PrimaryDirection   model.Direction
	SecondaryDirection model.Direction
	LastDirection      model.Direction
	IsRunning          bool
	OutgoingQueue      []DownstreamMessage
}

var SIDEBARS = []int{2423, 3917, 638, 3213, 1644, 5608, 1151,
	18128, 5065, 5715, 2449, 904, 147, 962}

func NewPlayer() *Player {
	spawn := &model.Position{
		X: 3200,
		Y: 3200,
	}
	player := &Player{
		Id:                 uuid.New(), // TODO: Load this from database or something
		Position:           spawn,
		LastKnownRegion:    spawn,
		SecondaryDirection: model.None,
		IsRunning:          true,
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
