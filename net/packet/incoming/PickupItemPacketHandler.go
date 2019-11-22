package incoming

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

type PickupItemPacketHandler struct{}

func (m *PickupItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	itemY := packet.ReadLEShort()
	itemId := packet.ReadShort()
	itemX := packet.ReadLEShort()

	m.pickupItemInternal(player, itemX, itemY, itemId)
}

func (m *PickupItemPacketHandler) pickupItemInternal(player *entity.Player, x, y, id uint16) {
	if player.DelayedDestination != nil {
		player.DelayedPacket = func() {
			m.pickupItemInternal(player, x, y, id)
		}
		return
	}

	groundItem := player.Region.FindGroundItemByPosition(int(id), &model.Position{X: x, Y: y})
	if groundItem == nil {
		return
	}

	player.Region.RemoveGroundItemIdAtPosition(int(id), &model.Position{X: x, Y: y})
	player.Inventory.AddItem(groundItem.ItemId, groundItem.Amount)
}
