package entity

import (
	"log"
	"rsps/model"
)

type MovementQueue struct {
	character *Character
	points    []*model.Point
}

func NewMovementQueue(c *Character) *MovementQueue {
	return &MovementQueue{
		character: c,
		points:    make([]*model.Point, 0),
	}
}

func (m *MovementQueue) Tick() {
	if len(m.points) > 0 {
		nextPos := m.points[0]
		m.points = m.points[1:]
		m.character.Position = nextPos.Position
		m.character.PrimaryDirection = nextPos.Direction
		//m.character.LastDirection = nextPos.Direction
		log.Printf("character %+v", m.character)
	}
}

func (m *MovementQueue) AddPosition(p *model.Position) {
	var last *model.Point
	if len(m.points) > 0 {
		last = m.points[len(m.points)-1]
	} else {
		last = &model.Point{
			Position:  m.character.Position,
			Direction: model.None,
		}
	}

	deltaX := int(p.X) - int(last.Position.X)
	deltaY := int(p.Y) - int(last.Position.Y)

	direction := model.DirectionFromDeltas(deltaX, deltaY)
	if direction != model.None {
		m.points = append(m.points, &model.Point{
			Position:  p,
			Direction: direction,
		})
	}
}
