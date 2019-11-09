package packet

import (
	"bytes"
	"rsps/model"
)

type PlayerUpdateType int

const (
	Idle PlayerUpdateType = iota
	Moving
	Running
	Teleport
)

type PlayerUpdatePacket struct {
	buf *bytes.Buffer

	updateRequired bool
	clearFlag      bool
	player         *model.Player
	otherPlayers   []interface{}
}

func NewPlayerUpdatePacket(player *model.Player) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{
		player: player,
		buf:    bytes.NewBuffer([]byte{81, 0, 0}),
	}
}

func (p *PlayerUpdatePacket) SetUpdateRequired(updateRequired bool) *PlayerUpdatePacket {
	p.updateRequired = updateRequired
	return p
}

func (p *PlayerUpdatePacket) SetClearFlag(clearFlag bool) *PlayerUpdatePacket {
	p.clearFlag = clearFlag
	return p
}

func (p *PlayerUpdatePacket) SetOtherPlayers(otherPlayers []interface{}) *PlayerUpdatePacket {
	p.otherPlayers = otherPlayers
	return p
}

func (p *PlayerUpdatePacket) Build() []byte {
	stream := model.NewStream()

	var updateType PlayerUpdateType
	if p.player.LastDirection != model.None {
		updateType = Running
	} else if p.player.PrimaryDirection != model.None {
		updateType = Moving
	} else {
		updateType = Idle
	}

	if p.updateRequired {
		stream.WriteBits(1, 1)
		switch updateType {
		case Idle:
			stream.WriteBits(2, 0)
			break
		case Moving:
			stream.WriteBits(2, 1)
			stream.WriteBits(3, uint(p.player.PrimaryDirection))
			stream.WriteBits(1, 1)
		case Running:
			stream.WriteBits(2, 2)
			stream.WriteBits(3, uint(p.player.PrimaryDirection))
			stream.WriteBits(3, uint(p.player.LastDirection))
			stream.WriteBits(1, 1)
			// TODO: Teleport
		}
	} else {
		stream.WriteBits(1, 0)
	}

	p.otherPlayers = make([]interface{}, 1)
	stream.WriteBits(8, uint(len(p.otherPlayers)-1))
	stream.WriteBits(11, 2047)
	p.buf.Write(stream.Flush())

	if p.updateRequired {
		updateMask := byte(0)

		updateMask |= 0x10 // appearance update
		// updateMask |= 4 // forced chat

		p.buf.WriteByte(updateMask)

		//p.buf.Write([]byte("Hello"))
		//p.buf.WriteByte(10)

		pa := &PlayerAppearance{
			Legs: 4730,
		}
		p.buf.Write(pa.ToBytes())
	}

	// calculate size of packet and set second word
	b := p.buf.Bytes()
	size := len(b) - 3
	b[1] = byte(size >> 8)
	b[2] = byte(size)

	return b
}
