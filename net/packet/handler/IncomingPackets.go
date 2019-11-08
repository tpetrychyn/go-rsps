package handler

import "rsps/net/packet"

const (
	GAME_MOVEMENT_OPCODE = 164
)

var IncomingPackets = map[byte]packet.PacketListener{
	GAME_MOVEMENT_OPCODE: new(MovementPacketHandler),
}
