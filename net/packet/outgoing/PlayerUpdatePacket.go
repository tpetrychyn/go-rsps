package outgoing

import (
	"bufio"
	"bytes"
	"rsps/model"
)

type PlayerUpdatePacket struct {
	player model.PlayerInterface
	buffer bytes.Buffer
}

func NewPlayerUpdatePacket(player model.PlayerInterface) *PlayerUpdatePacket {
	p := &PlayerUpdatePacket{
		player: player,
	}
	payload := p.Build()
	size := len(payload)
	p.buffer.WriteByte(81)
	p.buffer.WriteByte(byte(size >> 8))
	p.buffer.WriteByte(byte(size))
	p.buffer.Write(payload)

	return p
}

func (p *PlayerUpdatePacket) Write(writer *bufio.Writer) {
	writer.Write(p.buffer.Bytes())
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

	updateStream := model.NewStream()
	p.appendUpdates(updateStream, p.player, false)

	loadedPlayers := p.player.GetLoadedPlayers()
	stream.WriteBits(8, uint(len(loadedPlayers)))
	for _, v := range loadedPlayers {
		if !p.player.GetPosition().WithinRenderDistance(v.GetPosition()) || v.GetMarkedForDeletion() {
			stream.WriteBits(1, 1)
			stream.WriteBits(2, 3) // remove player from render?
			p.player.RemoveLoadedPlayer(v.GetId())
			continue
		}
		p.updateOtherPlayerMovement(stream, v)
		p.appendUpdates(updateStream, v, v.GetUpdateFlag().UpdateRequired)
	}

	localPlayers := p.player.GetNearbyPlayers()
localPlayerLoop:
	for _, v:= range localPlayers {
		if len(loadedPlayers) >= 79 {
			break
		}
		if v == p.player {
			continue
		}
		if !p.player.GetPosition().WithinRenderDistance(v.GetPosition()) {
			continue
		}
		for _, l := range loadedPlayers {
			if v.GetId() == l.GetId() {
				continue localPlayerLoop
			}
		}
		p.player.AddLoadedPlayer(v)
		p.addPlayer(stream, v)
		p.appendUpdates(updateStream, v, true)
	}

	updateBytes := updateStream.Flush()
	if len(updateBytes) > 1 {
		stream.WriteBits(11, 2047)
		buffer.Write(stream.Flush())
		buffer.Write(updateBytes)
	} else {
		buffer.Write(stream.Flush())
	}

	return buffer.Bytes()
}

func (p *PlayerUpdatePacket) appendUpdates(updateStream *model.Stream, target model.PlayerInterface, updateAppearance bool) {
	updateFlag := *target.GetUpdateFlag() // copy the targets updateFlag so we don't force the original players
	if updateAppearance {
		updateFlag.SetAppearance()
		updateFlag.SetFacePosition(updateFlag.FacePosition)
	}
	if updateFlag.UpdateRequired {
		var mask int
		// Setup the updateMask byte
		if updateFlag.ForcedMovement {
			mask |= 0x400
		}
		if updateFlag.Graphic {
			mask |= 0x100
		}
		if updateFlag.Animation {
			mask |= 0x8
		}
		if updateFlag.ForcedChat {
			mask |= 0x4
		}
		if updateFlag.Chat {
			mask |= 0x80
		} // TODO: Player ignore
		if updateFlag.EntityInteraction {
			mask |= 0x1
		}
		if updateFlag.Appearance {
			mask |= 0x10
		}
		if updateFlag.Face {
			mask |= 0x2
		}
		if updateFlag.SingleHit {
			mask |= 0x20
		}
		if updateFlag.DoubleHit {
			mask |= 0x200
		}

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
			p.updateAnimation(updateStream, target)
		}

		if updateFlag.ForcedChat {

		}

		if updateFlag.Chat {

		}

		if updateFlag.EntityInteraction {
			p.updateEntityInteraction(updateStream, target)
		}

		if updateFlag.Appearance {
			pa := &PlayerAppearance{Target: target, Equipment: target.GetEquipmentItemContainer()}
			updateStream.Write(pa.ToBytes())
		}

		if updateFlag.Face {
			p.facePosition(updateStream, target)
		}

		if updateFlag.SingleHit {
			p.updateSingleHit(updateStream)
		}

		if updateFlag.DoubleHit {
			p.updateDoubleHit(updateStream)
		}
	}
}

func (p *PlayerUpdatePacket) addPlayer(stream *model.Stream, player model.PlayerInterface) {
	stream.WriteBits(11, uint(player.GetId()))
	stream.WriteBits(1, 1)
	stream.WriteBits(1, 1)
	yDiff := int(player.GetPosition().Y) - int(p.player.GetPosition().Y)
	xDiff := int(player.GetPosition().X) - int(p.player.GetPosition().X)
	if xDiff < 0 {
		xDiff += 32
	} // 2^5 is 32, so xDiff needs to be between 0/32
	if yDiff < 0 {
		yDiff += 32
	}
	stream.WriteBits(5, uint(yDiff))
	stream.WriteBits(5, uint(xDiff))
}

func (p *PlayerUpdatePacket) updateOtherPlayerMovement(stream *model.Stream, target model.PlayerInterface) {
	if target.GetPrimaryDirection() == model.None {
		if target.GetUpdateFlag().UpdateRequired {
			stream.WriteBits(1, 1)
			stream.WriteBits(2, 0)
		} else {
			stream.WriteBits(1, 0)
		}
	} else if target.GetSecondaryDirection() == model.None {
		// walking
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 1)
		stream.WriteBits(3, uint(target.GetPrimaryDirection()))
		if target.GetUpdateFlag().UpdateRequired {
			stream.WriteBits(1, 1)
		} else {
			stream.WriteBits(1, 0)
		}
	} else {
		// Running
		stream.WriteBits(1, 1)
		stream.WriteBits(2, 2)
		stream.WriteBits(3, uint(target.GetPrimaryDirection()))
		stream.WriteBits(3, uint(target.GetSecondaryDirection()))
		if target.GetUpdateFlag().UpdateRequired {
			stream.WriteBits(1, 1)
		} else {
			stream.WriteBits(1, 0)
		}
	}
}

func (p *PlayerUpdatePacket) updateGraphics(stream *model.Stream) {
	stream.WriteWordLE(90)                            //graphicId
	stream.WriteInt((100 << 16) + (6553600 & 0xffff)) //height + delay
}

func (p *PlayerUpdatePacket) updateAnimation(stream *model.Stream, target model.PlayerInterface) {
	stream.WriteWordLE(uint(target.GetUpdateFlag().AnimationId)) //animId
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

func (p *PlayerUpdatePacket) facePosition(stream *model.Stream, target model.PlayerInterface) {
	var x uint16
	var y uint16
	if target.GetUpdateFlag().FacePosition == nil {
		x = 0
		y = 0
	} else {
		x = 2*target.GetUpdateFlag().FacePosition.X + 1
		y = 2*target.GetUpdateFlag().FacePosition.Y + 1
	}
	stream.WriteWordBEA(uint(x))
	stream.WriteWordLE(uint(y))
}

func (p *PlayerUpdatePacket) updateEntityInteraction(stream *model.Stream, target model.Character) {
	i := target.GetInteractingWith()
	if n, ok := i.(model.NpcInterface); ok {
		stream.WriteWordLE(uint(n.GetId()))
	} else if p, ok := i.(model.PlayerInterface); ok {
		stream.WriteWordLE(uint(32768 + p.GetId()))
	} else {
		stream.WriteWordLE(255)
	}
}
