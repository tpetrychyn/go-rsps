package incoming

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
)

type RemoveEquipmentPacketHandler struct {}

func (e *RemoveEquipmentPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	_ = packet.ReadLEShortA() //interfaceId?
	slot := packet.ReadShortA()
	id := packet.ReadShortA()

	if player.Equipment.Items[slot].ItemId != int(id) {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		player.Equipment.Items[slot] = &model.Item{}
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendItemContainerPacket{
			ItemContainer: player.Equipment,
			InterfaceId: model.EQUIPMENT_INTERFACE_ID,
		})
		return
	}

	player.Equipment.Items[slot] = &model.Item{}
	player.AddItem(int(id), 1)
	player.Equipment.SetItem(0, 0, int(slot))

	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: player.Equipment,
		InterfaceId: model.EQUIPMENT_INTERFACE_ID,
	})
}
