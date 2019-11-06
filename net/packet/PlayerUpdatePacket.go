package packet

import (
	"bytes"
)

type PlayerUpdateType int

const (
	Idle PlayerUpdateType = iota
	Moved
	Running
	Height
)

// [81 0 2 0 0] standing still
// [81 0 7 -90 1 -1 -64 1 0 0] move north
// [81 0 7 -70 1 -1 -64 1 0 0] move south

type PlayerUpdatePacket struct {
	buf *bytes.Buffer

	updateRequired bool
	typ            PlayerUpdateType
	clearFlag      bool
	otherPlayers   []interface{}
}

func NewPlayerUpdatePacket() *PlayerUpdatePacket {
	return &PlayerUpdatePacket{
		buf: bytes.NewBuffer([]byte{81, 0, 0}),
	}
}

func (p *PlayerUpdatePacket) SetUpdateRequired(updateRequired bool) *PlayerUpdatePacket {
	p.updateRequired = updateRequired

	return p
}

func (p *PlayerUpdatePacket) SetType(typ PlayerUpdateType) *PlayerUpdatePacket {
	p.typ = typ
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

var xlateDirectionToClient = []uint{ 1, 2, 4, 7, 6, 5, 3, 0 }
func (p *PlayerUpdatePacket) Build() []byte {
	stream := NewStream()

	if p.updateRequired {
		stream.WriteBits(1, 1)
		switch p.typ {
		case Idle:
			stream.WriteBits(2, 0)
			break
		case Moved:
			stream.WriteBits(2, 1)
			stream.WriteBits(3, xlateDirectionToClient[0])
			stream.WriteBits(1, 1)
		}
	} else {
		stream.WriteBits(1, 0)
	}

	//p.buf.Write(stream.Flush())

	p.otherPlayers = make([]interface{}, 1)
	stream.WriteBits(8, uint(len(p.otherPlayers)-1))
	stream.WriteBits(11, 2047)
	p.buf.Write(stream.Flush())

	if p.updateRequired {
		updateMask := byte(0)
		// forceChat
		//updateMask |= 4
		//p.buf.Write([]byte{updateMask})
		//p.buf.Write([]byte("null"))
		//p.buf.Write([]byte{10})

		updateMask |= 4 | 0x10
		p.buf.WriteByte(updateMask)

		p.buf.Write([]byte("hello"))
		p.buf.Write([]byte{10})

		p.buf.WriteByte(211) //this players update bit offset in the packet... wtf

		p.buf.WriteByte(0)    //player appearance 0
		p.buf.WriteByte(3)    // prayer icon
		p.buf.WriteByte(0xFF) //pk icon
		for i := 0; i < 12; i++ {
			p.buf.WriteByte(0)
		}
		for i := 0; i < 5; i++ {
			p.buf.WriteByte(0)
		}

		p.buf.Write([]byte{0x328 >> 8, 0x328 & 0xFF})
		p.buf.Write([]byte{0x337 >> 8, 0x337 & 0xFF})
		p.buf.Write([]byte{0x333 >> 8, 0x333 & 0xFF})
		p.buf.Write([]byte{0x334 >> 8, 0x334 & 0xFF})
		p.buf.Write([]byte{0x335 >> 8, 0x335 & 0xFF})
		p.buf.Write([]byte{0x336 >> 8, 0x336 & 0xFF})
		p.buf.Write([]byte{0x338 >> 8, 0x338 & 0xFF})

		//0 0 0 0 79 120 111 6
		p.buf.Write([]byte{0, 0, 0, 0, 79, 120, 111, 6}) //player name as int
		p.buf.WriteByte(3)
		p.buf.Write([]byte{0, 0})
	}
	// calculate size of packet and set second word
	b := p.buf.Bytes()
	size := len(b) - 3
	b[2] = byte(size)

	return b
}

//81 0 53 128 31 252 16 -48 0 -1 -1 0 0 0 0 0 0 1 26 0 1 0 0 0 1 10 0 0 0 0 0 3 40 3 55 3 51 3 52 3 53 3 54 3 56 0 0 0 0 79 120 111 6 3 0 0
//81 0 49 128 31 252 0 3 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3 40 3 55 3 51 3 52 3 53 3 54 3 56 0 0 0 0 0 0 0 0 100 10 10 0
