package model

type Movement struct {
	Position           *Position
	LastKnownRegion    *Position
	PrimaryDirection   Direction
	SecondaryDirection Direction
	LastDirection      Direction
	IsRunning          bool
	IsFrozen           bool
}

func (p *Movement) GetPosition() *Position {
	return p.Position
}

func (p *Movement) SetPosition(position *Position) {
	p.Position = position
}

func (p *Movement) GetLastKnownRegion() *Position {
	return p.LastKnownRegion
}

func (p *Movement) SetLastKnownRegion(position *Position) {
	p.LastKnownRegion = position
}

func (p *Movement) GetPrimaryDirection() Direction {
	return p.PrimaryDirection
}

func (p *Movement) SetPrimaryDirection(direction Direction) {
	p.PrimaryDirection = direction
}

func (p *Movement) GetSecondaryDirection() Direction {
	return p.SecondaryDirection
}

func (p *Movement) SetSecondaryDirection(direction Direction) {
	p.SecondaryDirection = direction
}

func (p *Movement) GetLastDirection() Direction {
	return p.LastDirection
}

func (p *Movement) SetLastDirection(direction Direction) {
	p.LastDirection = direction
}

func (p *Movement) GetIsFrozen() bool {
	return p.IsFrozen
}