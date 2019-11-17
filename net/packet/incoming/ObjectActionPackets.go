package incoming

import (
	"log"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

type ObjectActionPacket struct {}

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
	if player.Position.GetDistance(&model.Position{X: x, Y: y}) > 1 {
		player.DelayedPacket = func () {
			e.handleObjectActionOneInternal(player, x, y, id)
		}
		return
	}

	player.SkillHelper.AddExperience(model.Attack, 120000)
	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Your knowledge in Attack increases. Gained 200 experience in Attack."})
	log.Printf("Object Click1: x %+v, y %+v, id %+v", x, y, id)
}