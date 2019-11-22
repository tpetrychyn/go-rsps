package incoming

import (
	"encoding/binary"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

// TODO: Use model.point, it broke last attempt
type Point struct {
	X int8
	Y int8
}

type MovementPacketHandler struct {}

func (m *MovementPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	player.MovementQueue.Clear()
	player.UpdateFlag.SetEntityInteraction(nil)
	player.DelayedPacket = nil
	player.OngoingAction = nil

	if packet.Opcode == 248 {
		packet.Size -= 14
	}

	steps := int((packet.Size - 5) / 2)
	path := make([]Point, steps)
	firstStepX := packet.ReadLEShortA()
	for i := 0; i < steps; i++ {
		var point Point
		_ = binary.Read(packet.Buffer, binary.BigEndian, &point)
		path[i] = point
	}
	firstStepY := packet.ReadLEShort()

	positions := make([]*model.Position, steps+1)
	positions[0] = &model.Position{
		X: firstStepX,
		Y: firstStepY,
	}

	if player.GetPosition().GetDistance(positions[0]) >= 22 {
		return
	}

	for i := 0; i < steps; i++ {
		positions[i+1] = &model.Position{
			X: firstStepX + uint16(path[i].X),
			Y: firstStepY + uint16(path[i].Y),
			Z: 0,
		}
	}

	for _, v := range positions {
		player.MovementQueue.AddPosition(v)
	}

	if packet.Opcode == WALK_ON_COMMAND_OPCODE {
		player.DelayedDestination = positions[steps]
	}
}
