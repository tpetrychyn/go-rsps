package model

import (
	"math"
	"rsps/util"
)

type Position struct {
	X uint16
	Y uint16
	Z uint16
}

func (p *Position) AddX(n uint16) *Position {
	return &Position{
		X: p.X + 1,
		Y: p.Y,
		Z: 0,
	}
}

func (p *Position) AddY(n uint16) *Position {
	return &Position{
		X: p.X,
		Y: p.Y + 1,
		Z: 0,
	}
}

func (p *Position) GetRegionX() uint16 {
	return (p.X >> 3) - 6
}

func (p *Position) GetRegionY() uint16 {
	return (p.Y >> 3) - 6
}

func (p *Position) GetLocalX() uint16 {
	return p.X - 8 * p.GetRegionX()
}

func (p *Position) GetLocalY() uint16 {
	return p.Y - 8 * p.GetRegionY()
}

func (p *Position) isWithinDistance(other *Position, distance int) bool {
	deltaX := util.Abs(int(p.X) - int(other.X))
	deltaY := util.Abs(int(p.Y) - int(other.Y))
	return deltaX <= distance && deltaY <= distance
}

func (p *Position) GetDistance(other *Position) int {
	deltaX := p.X - other.X
	deltaY := p.Y - other.Y
	return int(math.Ceil(math.Sqrt(float64(deltaX * deltaX + deltaY * deltaY))))
}

func (p *Position) WithinRenderDistance(other *Position) bool {
	if p.Z != other.Z { return false }
	deltaX := int(other.X) - int(p.X)
	deltaY := int(other.Y) - int(p.Y)
	return deltaX <= 15 && deltaX >= -16 && deltaY <= 15 && deltaY >= -16
}