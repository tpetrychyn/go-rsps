package outgoing

import (
	"bufio"
	"bytes"
	"log"
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
		if !n.player.GetPosition().WithinRenderDistance(v.GetPosition()) {
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
	log.Printf("%+v", updateBytes)
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
	stream.WriteBits(14, uint(npc.GetNpcId()))
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
	stream.WriteBits(12, uint(npc.GetNpcType()))

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

	updateStream.WriteByte(byte(mask))

	if flag.Animation {
		n.updateAnimation(npc, updateStream)
	}
}

func (n *NpcUpdatePacket) updateAnimation(npc model.NpcInterface, updateStream *model.Stream) {
	anim := npc.GetUpdateFlag().AnimationId
	if anim < 0 {
		anim = 65535
	}
	updateStream.WriteWordLE(uint(anim)) //animId
	updateStream.WriteByte(1)
}