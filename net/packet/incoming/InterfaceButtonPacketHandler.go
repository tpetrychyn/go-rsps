package incoming

import (
	"log"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

const LOGOUT = 2458
const WALK_BUTTON = 152
const RUN_BUTTON = 153
const VARROCK_TELEPORT_BUTTON = 1164
const LUMBRIDGE_TELEPORT_BUTTON = 1167
const FALADOR_TELEPORT_BUTTON = 1170
const CAMELOT_TELEPORT_BUTTON = 1174
const ARDOUGNE_TELEPORT_BUTTON = 1540

type InterfaceButtonPacketHandler struct {}

func (i *InterfaceButtonPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	buttonId := packet.ReadShort()
	log.Printf("clicked button %+v", buttonId)

	switch buttonId {
	case LOGOUT:
		player.LogoutRequested = true
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.LogoutPacket{})

	case WALK_BUTTON:
		player.IsRunning = false
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.ConfigurationPacket{InterfaceId: 173,State:0})
	case RUN_BUTTON:
		player.IsRunning = true
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.ConfigurationPacket{InterfaceId: 173,State:1})

	case VARROCK_TELEPORT_BUTTON:
		player.Teleport(&model.Position{X: 3210, Y: 3424})
	case LUMBRIDGE_TELEPORT_BUTTON:
		player.Teleport(&model.Position{X: 3222, Y: 3218})
	case FALADOR_TELEPORT_BUTTON:
		player.Teleport(&model.Position{X: 2964, Y: 3378})
	case CAMELOT_TELEPORT_BUTTON:
		player.Teleport(&model.Position{X: 2757, Y: 3477})
	case ARDOUGNE_TELEPORT_BUTTON:
		player.Teleport(&model.Position{X: 2662, Y: 3305})
	}


}
