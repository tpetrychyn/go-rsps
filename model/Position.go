package model

type Position struct {
	X uint16
	Y uint16
	Z uint16
}

func NewPosition(x, y, z uint16) *Position {
	return &Position{
		X: x,
		Y: y,
		Z: z,
	}
}
