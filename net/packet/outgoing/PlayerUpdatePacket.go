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
	GetEquipmentItemContainer() *model.ItemContainer
	GetUpdateFlag() *model.UpdateFlag
}

type PlayerUpdatePacket struct {
	player       PlayerInterface
	otherPlayers []interface{}
}

func NewPlayerUpdatePacket(player PlayerInterface) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{
		player: player,
	}
}

func (p *PlayerUpdatePacket) SetOtherPlayers(otherPlayers []interface{}) *PlayerUpdatePacket {
	p.otherPlayers = otherPlayers
	return p
}

func (p *PlayerUpdatePacket) Write(writer *bufio.Writer) {
	payload := p.Build()
	size := len(payload)
	writer.WriteByte(81)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(payload)

	p.player.GetUpdateFlag().Clear()
}

func (p *PlayerUpdatePacket) Build() []byte {
	buffer := new(bytes.Buffer)
	stream := model.NewStream()

	updateFlag := p.player.GetUpdateFlag()

	if updateFlag.NeedsPlacement {
		// Teleported
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 3)
		stream.WriteBits(2, 0)
		stream.WriteBits(1, 1)
		if updateFlag.UpdateRequired {
			stream.WriteBits(1, 1)
		} else {
			stream.WriteBits(1, 0)
		}
		stream.WriteBits(7, uint(p.player.GetPosition().GetLocalY()))
		stream.WriteBits(7, uint(p.player.GetPosition().GetLocalX()))
	} else if p.player.GetSecondaryDirection() != model.None {
		// Running
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 2)
		stream.WriteBits(3, uint(p.player.GetPrimaryDirection()))
		stream.WriteBits(3, uint(p.player.GetSecondaryDirection()))
		if updateFlag.UpdateRequired {
			stream.WriteBits(1, 1)
		} else {
			stream.WriteBits(1, 0)
		}
	} else if p.player.GetPrimaryDirection() != model.None {
		// Walking
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 1)
		stream.WriteBits(3, uint(p.player.GetPrimaryDirection()))
		if updateFlag.UpdateRequired {
			stream.WriteBits(1, 1)
		} else {
			stream.WriteBits(1, 0)
		}
	} else if updateFlag.UpdateRequired {
		// Idle update
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 0)
	} else {
		// No update
		stream.WriteBits(1, 0)
	}

	p.otherPlayers = make([]interface{}, 1)
	stream.WriteBits(8, uint(len(p.otherPlayers)-1))

	var mask int
	updateStream := model.NewStream()
	if updateFlag.UpdateRequired {
		// Setup the updateMask byte
		if updateFlag.ForcedMovement { mask |= 0x400 }
		if updateFlag.Graphic { mask |= 0x100 }
		if updateFlag.Animation { mask |= 0x8 }
		if updateFlag.ForcedChat { mask |= 0x4 }
		if updateFlag.Chat { mask |= 0x80 } // TODO: Player ignore
		if updateFlag.EntityInteraction { mask |= 0x1 }
		if updateFlag.Appearance { mask |= 0x10 }
		if updateFlag.FacePosition { mask |= 0x2 }
		if updateFlag.SingleHit { mask |= 0x20 }
		if updateFlag.DoubleHit { mask |= 0x200 }

		if mask >= 0x100 {
			mask |= 0x40
			updateStream.WriteWordLE(uint(mask))
		} else {
			updateStream.WriteByte(byte(mask))
		}

		// Actual Flag updates begin here
		if updateFlag.ForcedMovement {

		}

		if updateFlag.Graphic {
			p.updateGraphics(updateStream)
		}

		if updateFlag.Animation {
			p.updateAnimation(updateStream)
		}

		if updateFlag.ForcedChat {

		}

		if updateFlag.Chat {

		}

		if updateFlag.EntityInteraction {

		}

		if updateFlag.Appearance {
			pa := &PlayerAppearance{Equipment: p.player.GetEquipmentItemContainer()}
			updateStream.Write(pa.ToBytes())
		}

		if updateFlag.FacePosition {

		}

		if updateFlag.SingleHit {
			p.updateSingleHit(updateStream)
		}

		if updateFlag.DoubleHit {
			p.updateDoubleHit(updateStream)
		}
	}

	if mask > 0 {
		stream.WriteBits(11, 2047)
		buffer.Write(stream.Flush())
		buffer.Write(updateStream.Flush())
	} else {
		buffer.Write(stream.Flush())
	}

	return buffer.Bytes()
}

func (p *PlayerUpdatePacket) updateGraphics(stream *model.Stream) {
	stream.WriteWordLE(90) //graphicId
	stream.WriteInt((100 << 16) + (6553600 & 0xffff)) //height + delay
}

func (p *PlayerUpdatePacket) updateAnimation(stream *model.Stream) {
	stream.WriteWordLE(865) //animId
	stream.WriteByte(^1 + 256)
}

func (p *PlayerUpdatePacket) updateSingleHit(stream *model.Stream) {
	stream.WriteByte(10)
	stream.WriteByte(1 + 128) // blue, red, green, yellow
	stream.WriteByte(^90 + 256)
	stream.WriteByte(99)
}

func (p *PlayerUpdatePacket) updateDoubleHit(stream *model.Stream) {
	stream.WriteByte(10)
	stream.WriteByte(128 - 1) // blue, red, green, yellow
	stream.WriteByte(90)
	stream.WriteByte(^99 + 256)
}