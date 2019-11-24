package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendPositionPacket struct {
	Position *model.Position
	Player   model.PlayerInterface
}

func (s *SendPositionPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(85)
	writer.WriteByte(-byte(s.Position.Y - 8*s.Player.GetLastKnownRegion().GetRegionY()))
	writer.WriteByte(-byte(s.Position.X - 8*s.Player.GetLastKnownRegion().GetRegionX()))
}
