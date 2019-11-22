package incoming

import (
	"log"
	"rsps/entity"
	"rsps/handler"
	"rsps/model"
	"rsps/net/packet"
)

type ObjectActionPacket struct{}

func (e *ObjectActionPacket) HandlePacket(player *entity.Player, packet *packet.Packet) {
	switch packet.Opcode {
	case OBJECT_ACTION_ONE_OPCODE:
		e.handleObjectActionOne(player, packet)
	}
}

func (e *ObjectActionPacket) handleObjectActionOne(player *entity.Player, packet *packet.Packet) {
	objectX := packet.ReadLEShortA()
	objectId := packet.ReadShort()
	objectY := packet.ReadShortA()

	e.handleObjectActionOneInternal(player, objectX, objectY, objectId)
}

func (e *ObjectActionPacket) handleObjectActionOneInternal(player *entity.Player, x, y, id uint16) {
	if player.DelayedDestination != nil {
		player.DelayedPacket = func() {
			e.handleObjectActionOneInternal(player, x, y, id)
		}
		return
	}

	player.UpdateFlag.SetFacePosition(&model.Position{X: x, Y: y})

	switch id {
	case 2213:
		player.Bank.OpenBank()
		return
	}

	if handler.ObjectIsWoodcuttingTree(int(id)) {
		handler.StartWoodcutting(int(id), &model.Position{X: x, Y: y}, player)
		return
	}

	log.Printf("Object Click1: x %+v, y %+v, id %+v", x, y, id)
}
