package incoming

import "rsps/net/packet"

const (
	EQUIP_ITEM_OPCODE = 41
	GAME_MOVEMENT_OPCODE = 164
	MINIMAP_MOVEMENT_OPCODE = 248
	INTERFACE_BUTTON_CLICK_OPCODE = 185
)

var Packets = map[byte]packet.PacketListener{
	GAME_MOVEMENT_OPCODE: new(MovementPacketHandler),
	MINIMAP_MOVEMENT_OPCODE: new(MovementPacketHandler),
	INTERFACE_BUTTON_CLICK_OPCODE: new(InterfaceButtonPacketHandler),
	EQUIP_ITEM_OPCODE: new(EquipItemPacketHandler),
}
