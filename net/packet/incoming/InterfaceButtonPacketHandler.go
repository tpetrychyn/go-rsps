package incoming

import (
	"log"
	"rsps/entity"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

const LOGOUT = 2458

type InterfaceButtonPacketHandler struct {}

func (i *InterfaceButtonPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	buttonId := packet.ReadShort()
	log.Printf("clicked button %+v", buttonId)

	switch buttonId {
	case LOGOUT:
		player.LogoutRequested = true
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.LogoutPacket{})
	}
}
