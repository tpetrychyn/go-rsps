package outgoing

import (
	"bufio"
	"rsps/model"
)

type RemoveGroundItemPacket struct {
	Player   model.PlayerInterface
	Position *model.Position
	ItemId   int
}

func (r *RemoveGroundItemPacket) Write(writer *bufio.Writer) {
	sendPositionPacket := &SendPositionPacket{Position: r.Position, Player: r.Player}
	sendPositionPacket.Write(writer)

	writer.WriteByte(156)
	writer.WriteByte(0 + 128)
	writer.WriteByte(byte(r.ItemId >> 8))
	writer.WriteByte(byte(r.ItemId))
}
