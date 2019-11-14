package incoming

import (
	"log"
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
	// TODO: x & y seem off by 1?
	if player.Position.GetDistance(&model.Position{X: x, Y: y}) > 0 {
		player.DelayedPacket = func() {
			m.pickupItemInternal(player, x, y, id)
		}
		return
	}
	groundItem := player.Region.FindGroundItemByPosition(int(id), &model.Position{
		X: x,
		Y: y,
	})
	// TODO: remove ground item
	if groundItem == nil {
		log.Printf("item %+v doesn't exist at x %+v, y %+v", id, x, y)
		return
	}
	player.Inventory.AddItem(groundItem.ItemId, groundItem.ItemAmount)
}
