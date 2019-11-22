package incoming

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

type CloseWindowPacketHandler struct {}

func (c *CloseWindowPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	// TODO: this is hardcoded as closing bank currently
	player.Inventory.CurrentInterfaceId = model.INVENTORY_INTERFACE_ID
	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: player.Inventory.ItemContainer,
		InterfaceId:   model.INVENTORY_INTERFACE_ID,
	})
}