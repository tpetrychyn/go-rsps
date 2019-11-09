package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Player struct {
	*model.Player
	*Character

	LoginState int
	//Outgoing []*packet.Packet
}

func NewPlayer() *Player {
	c := NewCharacter(&model.Position{
		X: 3222,
		Y: 3221,
		Z: 0,
	})
	return &Player{
		Character:  c,
	}
}

func (p *Player) PostUpdate() {
	p.Character.PostUpdate()
}

func (p *Player) Tick() {
	outgoing.SendMapRegion(&p.Player.Position)
	p.Character.Tick() // tick parent class
}
