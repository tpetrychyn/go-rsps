package outgoing

import (
	"bufio"
	"bytes"
	"rsps/model"
)

type PlayerInterface interface {
	GetLastDirection() model.Direction
	GetPrimaryDirection() model.Direction
	GetSecondaryDirection() model.Direction
	GetPosition() *model.Position
	GetLastKnownRegion() *model.Position
	GetEquipment() *model.ItemContainer
}

type PlayerUpdateType int

const (
	Idle PlayerUpdateType = iota
	Moving
	Running
	Teleport
)

type PlayerUpdatePacket struct {
	typ            PlayerUpdateType
	updateRequired bool
	clearFlag      bool
	player         PlayerInterface
	otherPlayers   []interface{}
}

func NewPlayerUpdatePacket(player PlayerInterface) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{
		player: player,
	}
}

func (p *PlayerUpdatePacket) SetUpdateRequired(updateRequired bool) *PlayerUpdatePacket {
	p.updateRequired = updateRequired
	return p
}

func (p *PlayerUpdatePacket) SetOtherPlayers(otherPlayers []interface{}) *PlayerUpdatePacket {
	p.otherPlayers = otherPlayers
	return p
}

func (p *PlayerUpdatePacket) SetTyp(typ PlayerUpdateType) *PlayerUpdatePacket {
	p.typ = typ
	return p
}

func (p *PlayerUpdatePacket) Write(writer *bufio.Writer) {
	payload := p.Build()
	size := len(payload)
	writer.WriteByte(81)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(payload)
}

func (p *PlayerUpdatePacket) Build() []byte {
	buffer := new(bytes.Buffer)
	stream := model.NewStream()

	var updateType = p.typ
	if p.typ == Teleport {
		updateType = Teleport
	} else if p.player.GetSecondaryDirection() != model.None {
		updateType = Running
	} else if p.player.GetPrimaryDirection() != model.None {
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
			stream.WriteBits(3, uint(p.player.GetPrimaryDirection()))
			stream.WriteBits(1, 1)
		case Running:
			stream.WriteBits(2, 2)
			stream.WriteBits(3, uint(p.player.GetPrimaryDirection()))
			stream.WriteBits(3, uint(p.player.GetSecondaryDirection()))
			stream.WriteBits(1, 1)
		case Teleport:
			stream.WriteBits(2, 3)
			stream.WriteBits(2, 0)
			stream.WriteBits(1, 1)
			stream.WriteBits(1, 1)
			stream.WriteBits(7, uint(p.player.GetPosition().GetLocalY()))
			stream.WriteBits(7, uint(p.player.GetPosition().GetLocalX()))
		}
	} else {
		stream.WriteBits(1, 0)
	}

	p.otherPlayers = make([]interface{}, 1)
	stream.WriteBits(8, uint(len(p.otherPlayers)-1))
	stream.WriteBits(11, 2047)
	buffer.Write(stream.Flush())

	if p.updateRequired {
		updateMask := byte(0)

		updateMask |= 0x10 // appearance update
		// updateMask |= 4 // forced chat

		buffer.WriteByte(updateMask)

		//p.buf.Write([]byte("Hello"))
		//p.buf.WriteByte(10)

		pa := &PlayerAppearance{
			Equipment: p.player.GetEquipment(),
		}
		buffer.Write(pa.ToBytes())
	}

	// calculate size of packet and set second word
	//b := p.buf.Bytes()
	//size := len(b) - 3
	//b[1] = byte(size >> 8)
	//b[2] = byte(size)

	return buffer.Bytes()
}
