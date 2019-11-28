package incoming

import "rsps/net/packet"

const (
	CHAT_OPCODE                   = 4
	SECOND_CLICK_NPC_OPCODE       = 17
	THIRD_CLICK_NPC_OPCODE        = 21
	CONTINUE_DIAGLOG_OPCODE       = 40
	EQUIP_ITEM_OPCODE             = 41
	BANK_TEN_OPCODE               = 43
	ATTACK_NPC_OPCODE             = 72
	DROP_ITEM_OPCODE              = 87
	WALK_ON_COMMAND_OPCODE        = 98
	COMMANDS_OPCODE               = 103
	BANK_FIVE_OPCODE              = 117
	CLICK_ITEM_OPCODE             = 122
	BANK_ALL_OPCODE               = 129
	CLOSE_WINDOW_OPCODE           = 130
	MAGE_NPC_OPCODE               = 131
	OBJECT_ACTION_ONE_OPCODE      = 132
	REMOVE_EQUIPMENT_OPCODE       = 145
	FIRST_CLICK_NPC_OPCODE        = 155
	GAME_MOVEMENT_OPCODE          = 164
	INTERFACE_BUTTON_CLICK_OPCODE = 185
	MOVE_ITEM_OPCODE              = 214
	PICKUP_ITEM_OPCODE            = 236
	MINIMAP_MOVEMENT_OPCODE       = 248
)

var Packets = map[byte]packet.PacketListener{
	CHAT_OPCODE:                   new(ChatPacketHandler),
	ATTACK_NPC_OPCODE:             new(ClickNpcPacketHandler),
	SECOND_CLICK_NPC_OPCODE:       new(ClickNpcPacketHandler),
	CONTINUE_DIAGLOG_OPCODE:       new(ContinueDialogPacketHandler),
	COMMANDS_OPCODE:               new(CommandsPacketHandler),
	CLICK_ITEM_OPCODE:             new(ClickItemPacketHandler),
	BANK_FIVE_OPCODE:              new(BankFivePacketHandler),
	BANK_TEN_OPCODE:               new(BankTenPacketHandler),
	BANK_ALL_OPCODE:               new(BankAllPacketHandler),
	CLOSE_WINDOW_OPCODE:           new(CloseWindowPacketHandler),
	DROP_ITEM_OPCODE:              new(DropItemPacketHandler),
	EQUIP_ITEM_OPCODE:             new(EquipItemPacketHandler),
	WALK_ON_COMMAND_OPCODE:        new(MovementPacketHandler),
	GAME_MOVEMENT_OPCODE:          new(MovementPacketHandler),
	MINIMAP_MOVEMENT_OPCODE:       new(MovementPacketHandler),
	INTERFACE_BUTTON_CLICK_OPCODE: new(InterfaceButtonPacketHandler),
	MOVE_ITEM_OPCODE:              new(MoveItemPacketHandler),
	REMOVE_EQUIPMENT_OPCODE:       new(RemoveItemPacketHandler),
	PICKUP_ITEM_OPCODE:            new(PickupItemPacketHandler),
	OBJECT_ACTION_ONE_OPCODE:      new(ObjectActionPacket),
}
