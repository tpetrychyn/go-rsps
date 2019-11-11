package outgoing

import (
	"bufio"
	"rsps/model"
)

type CreateGroundItemPacket struct {
	Position   *model.Position
	Player     PlayerInterface
	ItemId     int
	ItemAmount int
}

func (c *CreateGroundItemPacket) Write(writer *bufio.Writer) {
	sendPositionPacket := &SendPositionPacket{Position: c.Position, Player: c.Player}
	sendPositionPacket.Write(writer)

	writer.WriteByte(44)
	writer.WriteByte(byte(c.ItemId + 128))
	writer.WriteByte(byte(c.ItemId >> 8))
	writer.WriteByte(byte(c.ItemAmount >> 8))
	writer.WriteByte(byte(c.ItemAmount))
	writer.WriteByte(0)
}
