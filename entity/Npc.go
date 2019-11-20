package entity

import (
	"github.com/google/uuid"
	"log"
	"rsps/model"
)

type Npc struct {
	*model.Movement
	Id         uuid.UUID
	NpcId      int
	NpcType    int
	UpdateFlag *model.UpdateFlag
}

func NewNpc() *Npc {
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
		Id:    uuid.UUID{},
		UpdateFlag: &model.UpdateFlag{},
		NpcId: 1,
		NpcType: 2,
	}
}

var flip = 0
func (n *Npc) Tick() {
	//if n.UpdateFlag.AnimationDuration <= 0 {
		if flip == 0 {
			n.UpdateFlag.SetAnimation(422, 2)
			flip = 1
		} else if flip == 1 {
			//n.UpdateFlag.ClearAnimation()
			flip = 2
		} else if flip == 2 {
			n.UpdateFlag.SetAnimation(404, 2)
			log.Printf("set to 404")
			flip = 3
		} else {
			//n.UpdateFlag.ClearAnimation()
			flip = 0
		}
	//}
}

func (n *Npc) PostUpdate() {
	n.UpdateFlag.Clear()
}

func (n *Npc) GetId() uuid.UUID {
	return n.Id
}
func (n *Npc) GetNpcId() int {
	return n.NpcId
}
func (n *Npc) GetNpcType() int {
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
