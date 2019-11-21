package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendSongPacket struct {
	Song int
}

func (s *SendSongPacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(74)
	stream.WriteWord(uint(s.Song))
	writer.Write(stream.Flush())
}

type SendSoundPacket struct {
	Sound int
	Volume int
	Delay int
}

func (s *SendSoundPacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(174)
	stream.WriteWord(uint(s.Sound))
	stream.WriteByte(byte(s.Volume))
	stream.WriteWord(uint(s.Delay))
	writer.Write(stream.Flush())
}
