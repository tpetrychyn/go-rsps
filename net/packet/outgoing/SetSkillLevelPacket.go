package outgoing

import (
	"bufio"
	"rsps/model"
)

type SetSkillLevelPacket struct {
	SkillNum   int
	Level      int
	Experience int
}

func (s *SetSkillLevelPacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(134)
	stream.WriteByte(byte(s.SkillNum))
	stream.WriteDWord_v1(s.Experience)
	stream.WriteByte(byte(s.Level))

	writer.Write(stream.Flush())
}
