package model

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
