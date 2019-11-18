package incoming

import (
	"rsps/entity"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

type ContinueDialogPacketHandler struct {}

func (c *ContinueDialogPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.ClearInterfacePacket{})
}
