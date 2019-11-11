package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/util"
)

type MovementQueue struct {
	character Character
	points    []*model.Point
}

func NewMovementQueue(c Character) *MovementQueue {
	return &MovementQueue{
		character: c,
		points:    make([]*model.Point, 0),
	}
}

func (m *MovementQueue) Tick() {
	if len(m.points) > 0 {
		lastPosition := m.character.GetPosition()
		walkPoint := m.points[0]
		m.points = m.points[1:]
		var runPoint *model.Point
		if p, ok := m.character.(*Player); ok && p.IsRunning && len(m.points) > 0 {
			runPoint = m.points[0]
			m.points = m.points[1:]
		}

		if walkPoint != nil && walkPoint.Direction != model.None {
			m.character.SetPosition(walkPoint.Position)
			m.character.SetPrimaryDirection(walkPoint.Direction)
			m.character.SetLastDirection(walkPoint.Direction)
		}

		if runPoint != nil && runPoint.Direction != model.None {
			m.character.SetPosition(runPoint.Position)
			m.character.SetSecondaryDirection(runPoint.Direction)
			m.character.SetLastDirection(runPoint.Direction)
		}

		diffX := m.character.GetPosition().X - m.character.GetLastKnownRegion().GetRegionX()*8
		diffY := m.character.GetPosition().Y - m.character.GetLastKnownRegion().GetRegionY()*8
		if p, ok := m.character.(*Player); ok {
			if diffX < 16 || diffX >= 88 || diffY < 16 || diffY >= 88 {
				p.LastKnownRegion = p.Position
				p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.MapRegionPacket{Position: p.Position})
			}

			if GetRegionIdByPosition(p.Position) != GetRegionIdByPosition(lastPosition) {
				WorldProvider().AddPlayerToRegion(p)
			}
		}
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
			Position:  m.character.GetPosition(),
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

	max := util.Abs(deltaX)
	if util.Abs(deltaY) > util.Abs(deltaX) {
		max = util.Abs(deltaY)
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
