package incoming

import "rsps/net/packet"

const (
	CHAT_OPCODE = 4
	EQUIP_ITEM_OPCODE = 41
	REMOVE_EQUIPMENT_OPCODE = 145
	GAME_MOVEMENT_OPCODE = 164
	INTERFACE_BUTTON_CLICK_OPCODE = 185
	MOVE_ITEM_OPCODE = 214
	MINIMAP_MOVEMENT_OPCODE = 248
)

var Packets = map[byte]packet.PacketListener{
	CHAT_OPCODE: new(ChatPacketHandler),
	GAME_MOVEMENT_OPCODE: new(MovementPacketHandler),
	MINIMAP_MOVEMENT_OPCODE: new(MovementPacketHandler),
	INTERFACE_BUTTON_CLICK_OPCODE: new(InterfaceButtonPacketHandler),
	EQUIP_ITEM_OPCODE: new(EquipItemPacketHandler),
	MOVE_ITEM_OPCODE: new(MoveItemPacketHandler),
	REMOVE_EQUIPMENT_OPCODE: new(RemoveEquipmentPacketHandler),
}
