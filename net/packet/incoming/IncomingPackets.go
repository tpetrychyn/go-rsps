package incoming

import "rsps/net/packet"

const (
	CHAT_OPCODE                   = 4
	EQUIP_ITEM_OPCODE             = 41
	WALK_ON_COMMAND_OPCODE        = 98
	OBJECT_ACTION_ONE_OPCODE      = 132
	REMOVE_EQUIPMENT_OPCODE       = 145
	GAME_MOVEMENT_OPCODE          = 164
	INTERFACE_BUTTON_CLICK_OPCODE = 185
	MOVE_ITEM_OPCODE              = 214
	PICKUP_ITEM_OPCODE            = 236
	MINIMAP_MOVEMENT_OPCODE       = 248
)

var Packets = map[byte]packet.PacketListener{
	CHAT_OPCODE:                   new(ChatPacketHandler),
	EQUIP_ITEM_OPCODE:             new(EquipItemPacketHandler),
	WALK_ON_COMMAND_OPCODE:        new(MovementPacketHandler),
	GAME_MOVEMENT_OPCODE:          new(MovementPacketHandler),
	MINIMAP_MOVEMENT_OPCODE:       new(MovementPacketHandler),
	INTERFACE_BUTTON_CLICK_OPCODE: new(InterfaceButtonPacketHandler),
	MOVE_ITEM_OPCODE:              new(MoveItemPacketHandler),
	REMOVE_EQUIPMENT_OPCODE:       new(RemoveEquipmentPacketHandler),
	PICKUP_ITEM_OPCODE:            new(PickupItemPacketHandler),
	OBJECT_ACTION_ONE_OPCODE:      new(ObjectActionPacket),
}
