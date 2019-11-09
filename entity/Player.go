package entity

import (
	"rsps/model"
)


type Player struct {
	*Character
}

func NewPlayer() *Player {
	c := NewCharacter(&model.Position{
		X: 3222,
		Y: 3221,
		Z: 0,
	})
	p := &Player{
		Character: c,
	}
	//go p.Tick()

	return p
}

func (p *Player) PostUpdate() {
	p.Character.PostUpdate()
}

func (p *Player) Tick() {
	//for {
	p.Character.Tick() // tick parent class
	//
	//	time.Sleep(600 * time.Millisecond)
	//
	//	p.PostUpdate()
	//}

}
