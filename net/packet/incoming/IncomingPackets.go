package incoming

import "rsps/net/packet"

const (
	GAME_MOVEMENT_OPCODE = 164
	MINIMAP_MOVEMENT_OPCODE = 248
)

var Packets = map[byte]packet.PacketListener{
	GAME_MOVEMENT_OPCODE: new(MovementPacketHandler),
	MINIMAP_MOVEMENT_OPCODE: new(MovementPacketHandler),
}
