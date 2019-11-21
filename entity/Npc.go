package entity

import (
	"math/rand"
	"rsps/model"
	"time"
)

type Npc struct {
	*model.Movement
	MovementQueue     *MovementQueue
	Id                int
	NpcType           int
	CurrentHitpoints  int
	MaxHitpoints      int
	UpdateFlag        *model.UpdateFlag
	Killer            *Player
	MarkedForDeletion bool
}

func NewNpc(id int) *Npc {
	spawn := &model.Position{
		X: 3201,
		Y: 3200,
	}
	npc := &Npc{
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			PrimaryDirection:   model.None,
			SecondaryDirection: model.None,
		},
		Id:               id,
		NpcType:          2,
		CurrentHitpoints: 10,
		MaxHitpoints:     10,
		UpdateFlag:       &model.UpdateFlag{},
	}
	mq := NewMovementQueue(npc)
	npc.MovementQueue = mq
	return npc
}

var spawn = &model.Position{
	X: 3200,
	Y: 3200,
	Z: 0,
}

func (n *Npc) Tick() {
	n.MovementQueue.Clear()
	// TODO: refactor all this lol
	if len(n.MovementQueue.points) == 0 && n.UpdateFlag.InteractingWith == nil {
		n.MovementQueue.AddPosition(&model.Position{
			X: spawn.X + uint16(rand.Intn(3 - 0)),
			Y: spawn.Y + uint16(rand.Intn(3 - 0)),
			Z: 0,
		})
	}

	if n.UpdateFlag.InteractingWith != nil && n.Position.GetDistance(n.UpdateFlag.InteractingWith.GetPosition()) > 1 {
		n.MovementQueue.AddPosition(n.UpdateFlag.InteractingWith.GetPosition())
	}

	if n.UpdateFlag.InteractingWith != nil && n.Position.GetDistance(n.UpdateFlag.InteractingWith.GetPosition()) == 0 {
		n.MovementQueue.AddPosition(n.UpdateFlag.InteractingWith.GetPosition().AddX(1))
	}

	n.MovementQueue.Tick()
	if n.CurrentHitpoints <= 0 && !n.MarkedForDeletion {
		n.UpdateFlag.SetAnimation(836, 2)
		n.MarkedForDeletion = true
		// TODO: stop this firing twice but also dont immediately remove from players ^
		go func() {
			<-time.After(1 * time.Second)
			world := WorldProvider()
			region := world.GetRegion(GetRegionIdByPosition(n.Position))
			if n.Killer != nil {
				region.CreateGroundItemAtPosition(n.Killer, &model.Item{
					ItemId: 995,
					Amount: 10000,
				}, n.Position)
			}
			world.RemoveNpc(n.Id)
		}()
	}
}

func (n *Npc) PostUpdate() {
	n.PrimaryDirection = model.None
	n.SecondaryDirection = model.None
	n.LastDirection = model.None
	n.UpdateFlag.Clear()
}

func (n *Npc) GetId() int {
	return n.Id
}
func (n *Npc) GetType() int {
	return n.NpcType
}
func (n *Npc) GetPrimaryDirection() model.Direction {
	return n.PrimaryDirection
}
func (n *Npc) GetUpdateFlag() *model.UpdateFlag {
	return n.UpdateFlag
}

func (n *Npc) GetPosition() *model.Position {
	return n.Position
}

func (n *Npc) GetInteractingWith() model.Character {
	return n.UpdateFlag.InteractingWith
}

func (n *Npc) GetCurrentHitpoints() int {
	return n.CurrentHitpoints
}

func (n *Npc) GetMaxHitpoints() int {
	return n.MaxHitpoints
}

func (n *Npc) TakeDamage(damage int) {
	n.UpdateFlag.SetSingleHit(damage)
	n.CurrentHitpoints -= damage
}

func (n *Npc) GetMarkedForDeletion() bool {
	return n.MarkedForDeletion
}
