package entity

import (
	"rsps/model"
	"time"
)

type Npc struct {
	*model.Movement
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

	return &Npc{
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
}

func (n *Npc) Tick() {
	if n.CurrentHitpoints <= 0 && !n.MarkedForDeletion {
		n.UpdateFlag.SetAnimation(836, 2)
		n.MarkedForDeletion = true
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
