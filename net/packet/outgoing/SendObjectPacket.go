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
	if a.ObjectId == -1 {
		stream.WriteByte(101)
		stream.WriteByte(-byte(a.Typ<<2 + (a.Face & 3)))
		stream.WriteByte(0)
	} else {
		stream.WriteByte(151)
		stream.WriteByte(128 + byte(a.Position.Z))
		stream.WriteWordLE(uint(a.ObjectId))
		stream.WriteByte(128 - byte((a.Typ << 2) + (a.Face & 3)))
	}

	writer.Write(stream.Flush())
}
