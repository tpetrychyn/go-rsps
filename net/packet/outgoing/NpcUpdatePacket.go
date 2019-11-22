package outgoing

import (
	"bufio"
	"bytes"
	"rsps/model"
)

type NpcUpdatePacket struct {
	player model.PlayerInterface
	buffer bytes.Buffer
}

func NewNpcUpdatePacket(player model.PlayerInterface) *NpcUpdatePacket {
	n := &NpcUpdatePacket{
		player: player,
	}

	payload := n.Build()
	size := len(payload)
	n.buffer.WriteByte(65)
	n.buffer.WriteByte(byte(size >> 8))
	n.buffer.WriteByte(byte(size))
	n.buffer.Write(payload)

	return n
}

func (n *NpcUpdatePacket) Write(writer *bufio.Writer) {
	writer.Write(n.buffer.Bytes())
}

func (n *NpcUpdatePacket) Build() []byte {
	buffer := new(bytes.Buffer)
	stream := model.NewStream()
	updateStream := model.NewStream()

	loadedNpcs := n.player.GetLoadedNpcs()
	stream.WriteBits(8, uint(len(loadedNpcs)))
	for _, v := range loadedNpcs {
		if !n.player.GetPosition().WithinRenderDistance(v.GetPosition()) || v.GetMarkedForDeletion() {
			stream.WriteBits(1, 1)
			stream.WriteBits(2, 3)
			n.player.RemoveLoadedNpc(v.GetId())
			continue
		}
		n.updateNpcMovement(stream, v)
		n.appendUpdates(v, updateStream)
	}

	localNpcs := n.player.GetNearbyNpcs()
localNpcsLoop:
	for _, v := range localNpcs {
		if len(loadedNpcs) >= 79 {
			break
		}
		if !n.player.GetPosition().WithinRenderDistance(v.GetPosition()) {
			continue
		}
		for _, l := range loadedNpcs {
			if v.GetId() == l.GetId() {
				continue localNpcsLoop
			}
		}

		n.player.AddLoadedNpc(v)
		n.addNewNpc(v, stream)
		n.appendUpdates(v, updateStream)
	}

	updateBytes := updateStream.Flush()
	if len(updateBytes) > 1 {
		stream.WriteBits(14, 16383)
		buffer.Write(stream.Flush())
		buffer.Write(updateBytes)
	} else {
		buffer.Write(stream.Flush())
	}

	return buffer.Bytes()
}

func (n *NpcUpdatePacket) updateNpcMovement(stream *model.Stream, target model.NpcInterface) {
	if target.GetSecondaryDirection() == model.None {
		if target.GetPrimaryDirection() == model.None {
			if target.GetUpdateFlag().UpdateRequired {
				stream.WriteBits(1, 1)
				stream.WriteBits(2, 0)
			} else {
				stream.WriteBits(1, 0)
			}
		} else {
			stream.WriteBits(1, 1)
			stream.WriteBits(2, 1)
			stream.WriteBits(3, uint(target.GetPrimaryDirection()))
			if target.GetUpdateFlag().UpdateRequired {
				stream.WriteBits(1, 1)
			} else {
				stream.WriteBits(1, 0)
			}
		}
	} else {
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

func (n *NpcUpdatePacket) addNewNpc(npc model.NpcInterface, stream *model.Stream) {
	stream.WriteBits(14, uint(npc.GetId()))
	z := int(npc.GetPosition().Y) - int(n.player.GetPosition().Y)
	if z < 0 {
		z += 32
	}
	stream.WriteBits(5, uint(z))
	z = int(npc.GetPosition().X) - int(n.player.GetPosition().X)
	if z < 0 {
		z += 32
	}
	stream.WriteBits(5, uint(z))

	stream.WriteBits(1, 0)
	stream.WriteBits(12, uint(npc.GetType()))

	if npc.GetUpdateFlag().UpdateRequired {
		stream.WriteBits(1, 1)
	} else {
		stream.WriteBits(1, 0)
	}
}

func (n *NpcUpdatePacket) appendUpdates(npc model.NpcInterface, updateStream *model.Stream) {
	var mask = 0
	flag := *npc.GetUpdateFlag()

	if flag.Animation {
		mask |= 0x10
	}
	if flag.SingleHit {
		mask |= 0x08
	}
	if flag.Graphic {
		mask |= 0x8
	}
	if flag.EntityInteraction {
		mask |= 0x20
	}
	if flag.ForcedChat != "" {
		mask |= 0x1
	}
	if flag.DoubleHit {
		mask |= 0x40
	}
	if flag.Transform {
		mask |= 0x2
	}
	if flag.Face {
		mask |= 0x4
	}

	updateStream.WriteByte(byte(mask))

	if flag.Animation {
		n.updateAnimation(updateStream, npc)
	}

	if flag.SingleHit {
		n.updateSingleHit(updateStream, npc)
	}

	if flag.Graphic {
	}

	if flag.EntityInteraction {
		n.updateEntityInteraction(updateStream, npc)
	}

	if flag.ForcedChat != "" {
		n.updateForcedChat(updateStream, npc)
	}

	if flag.Face {
		n.updateFace(updateStream, npc)
	}
}

func (n *NpcUpdatePacket) updateAnimation(updateStream *model.Stream, npc model.NpcInterface) {
	anim := npc.GetUpdateFlag().AnimationId
	if anim < 0 {
		anim = 65535
	}
	updateStream.WriteWordLE(uint(anim)) //animId
	updateStream.WriteByte(1)
}

func (n *NpcUpdatePacket) updateEntityInteraction(stream *model.Stream, npc model.Character) {
	i := npc.GetInteractingWith()
	if n, ok := i.(model.NpcInterface); ok {
		stream.WriteWordLE(uint(n.GetId()))
	} else if p, ok := i.(model.PlayerInterface); ok {
		stream.WriteWord(uint(32768 + p.GetId()))
	} else {
		stream.WriteWordLE(255)
	}
}

func (n *NpcUpdatePacket) updateSingleHit(stream *model.Stream, npc model.Character) {
	damage := npc.GetUpdateFlag().SingleHitDamage
	stream.WriteByte(byte(damage) + 128)
	if damage > 0 {
		stream.WriteByte(255)
	} else {
		stream.WriteByte(0)
	}
	stream.WriteByte(byte(npc.GetCurrentHitpoints()) + 128)
	stream.WriteByte(byte(npc.GetMaxHitpoints()))
}

func (n *NpcUpdatePacket) updateFace(updateStream *model.Stream, npc model.NpcInterface) {
	var x uint16
	var y uint16
	if npc.GetUpdateFlag().FacePosition == nil {
		x = 0
		y = 0
	} else {
		x = 2*npc.GetUpdateFlag().FacePosition.X + 1
		y = 2*npc.GetUpdateFlag().FacePosition.Y + 1
	}
	updateStream.WriteWordLE(uint(x))
	updateStream.WriteWordLE(uint(y))
}

func (n *NpcUpdatePacket) updateForcedChat(updateStream *model.Stream, npc model.NpcInterface) {
	updateStream.Write([]byte(npc.GetUpdateFlag().ForcedChat))
	updateStream.WriteByte(10)
}
