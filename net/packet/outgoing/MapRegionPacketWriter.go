package outgoing

import (
	"bufio"
	"rsps/model"
)

type MapRegionPacket struct {
	Position *model.Position
}

func (m *MapRegionPacket) Write(writer *bufio.Writer) {
	s := model.NewStream()
	s.WriteWordA(uint(m.Position.GetRegionX() + 6))
	s.WriteWord(uint(m.Position.GetRegionY() + 6))

	writer.WriteByte(73)
	writer.Write(s.Flush())
}

