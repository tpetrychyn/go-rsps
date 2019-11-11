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
