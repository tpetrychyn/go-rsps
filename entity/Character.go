package entity

import (
	"rsps/model"
)

type Character struct {
	MovementQueue    *MovementQueue
	Position         *model.Position
	PrimaryDirection model.Direction
	LastDirection    model.Direction
}

func NewCharacter(p *model.Position) *Character {
	c := &Character{
		Position:         p,
		PrimaryDirection: model.None,
		LastDirection:    model.None,
	}
	c.MovementQueue = NewMovementQueue(c)

	return c
}

func (c *Character) Tick() {
	c.MovementQueue.Tick()
}

func (c *Character) PostUpdate() {
	c.PrimaryDirection = model.None
	c.LastDirection = model.None
}
