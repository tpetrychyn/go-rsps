package entity

import (
	"log"
	"rsps/model"
)

type Character struct {
	Position         *model.Position
	PrimaryDirection model.Direction
	LastDirection    model.Direction
	MovementQueue    *MovementQueue
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
	log.Printf("c character %+v", c)
}

func (c *Character) PostUpdate() {
	c.PrimaryDirection = model.None
	c.LastDirection = model.None
}
