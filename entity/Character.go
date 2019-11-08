package entity

import (
	"log"
	"rsps/model"
)

type Character struct {
	Position      model.Position
	MovementQueue []*model.Position
}

func (c *Character) Tick() {
	if len(c.MovementQueue) > 0 {
		nextPos := c.MovementQueue[0]
		c.MovementQueue = c.MovementQueue[1:]
		log.Printf("%+v", nextPos)
	}
}
