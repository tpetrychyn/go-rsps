package entity

import (
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
	Killer            model.Character
	MarkedForDeletion bool
	IsDying           bool

	GlobalTickCount int
	AttackSpeed     int
	OngoingAction   model.OngoingAction
}

func NewNpc(id int, npcType int, position *model.Position) *Npc {
	// TODO: npcs will be gone when region is deleted, need to respawn them in constructor probably
	attackSpeed := 5
	npc := &Npc{
		Movement: &model.Movement{
			Position:           position,
			LastKnownRegion:    position,
			PrimaryDirection:   model.None,
			SecondaryDirection: model.None,
		},
		Id:               id,
		NpcType:          npcType,
		CurrentHitpoints: 10,
		MaxHitpoints:     10,
		AttackSpeed:      attackSpeed,
		UpdateFlag:       &model.UpdateFlag{},
	}
	mq := NewMovementQueue(npc)
	npc.MovementQueue = mq
	return npc
}

func (n *Npc) Tick() {
	n.MovementQueue.Tick()

	if n.GlobalTickCount > 0 {
		n.GlobalTickCount--
	}

	if n.OngoingAction != nil {
		n.OngoingAction.Tick()
	}

	if n.CurrentHitpoints <= 0 && !n.IsDying {
		n.UpdateFlag.SetAnimation(836, 2)
		n.IsDying = true
		n.IsFrozen = true
		go func() {
			<-time.After(1 * time.Second)
			n.MarkedForDeletion = true
			world := WorldProvider()
			region := world.GetRegion(GetRegionIdByPosition(n.Position))
			if n.Killer != nil {
				if p, ok := n.Killer.(*Player); ok {
					region.CreateGroundItemAtPosition(p, &model.Item{
						ItemId: 995,
						Amount: 10000,
					}, n.Position)
				}
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

func (n *Npc) GetGlobalTickCount() int {
	return n.GlobalTickCount
}

func (n *Npc) SetGlobalTickCount(g int) {
	n.GlobalTickCount = g
}

func (n *Npc) GetOngoingAction() model.OngoingAction {
	return n.OngoingAction
}

func (n *Npc) SetOngoingAction(action model.OngoingAction) {
	n.OngoingAction = action
}
