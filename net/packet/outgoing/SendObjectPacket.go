package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendObjectPacket struct {
	ObjectId int
	Position *model.Position
	Face     int
	Typ      int

	Player model.PlayerInterface
}

func (a *SendObjectPacket) Write(writer *bufio.Writer) {
	sendPositionPacket := &SendPositionPacket{Position: a.Position, Player: a.Player}
	sendPositionPacket.Write(writer)

	stream := model.NewStream()
	stream.WriteByte(151)
	stream.WriteByte(byte(128 - a.Position.Z))
	stream.WriteWordLE(uint(a.ObjectId))
	stream.WriteByte(byte(128 - (a.Typ << 2) + (a.Face & 3)))

	writer.Write(stream.Flush())
}
