package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendChatInterfacePacket struct {
	InterfaceId uint
}

func (s *SendChatInterfacePacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(164)
	stream.WriteWordLE(s.InterfaceId)

	writer.Write(stream.Flush())
}
