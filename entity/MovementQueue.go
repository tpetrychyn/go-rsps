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
		log.Printf("currentPos %+v, target %+v", m.character.Position, nextPos.Position)
	}
}

func (m *MovementQueue) Reset() {
	m.points = make([]*model.Point, 0)
}

func (m *MovementQueue) getLast() *model.Point {
	var last *model.Point
	if len(m.points) > 0 {
		last = m.points[len(m.points)-1]
	} else {
		last = &model.Point{
			Position:  m.character.Position,
			Direction: model.None,
		}
	}

	return last
}

func (m *MovementQueue) AddPosition(p *model.Position) {
	last := m.getLast()
	x := int(p.X)
	y := int(p.Y)

	deltaX := x - int(last.Position.X)
	deltaY := y - int(last.Position.Y)

	max := Abs(deltaX)
	if Abs(deltaY) > Abs(deltaX) {
		max = Abs(deltaY)
	}

	for i := 0; i < max; i++ {
		if deltaX < 0 {
			deltaX++
		} else if deltaX > 0 {
			deltaX--
		}
		if deltaY < 0 {
			deltaY++
		} else if deltaY > 0 {
			deltaY--
		}

		m.addStep(x-deltaX, y-deltaY, 0)
	}
}

func (m *MovementQueue) addStep(x, y, z int) {
	last := m.getLast()
	deltaX := x - int(last.Position.X)
	deltaY := y - int(last.Position.Y)
	direction := model.DirectionFromDeltas(deltaX, deltaY)
	if direction != model.None {
		m.points = append(m.points, &model.Point{
			Position: &model.Position{
				X: uint16(x),
				Y: uint16(y),
				Z: 0,
			},
			Direction: direction,
		})
	}
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
