package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendTextInterfacePacket struct {
	InterfaceId uint
	Message     string
}

func (s *SendTextInterfacePacket) Write(writer *bufio.Writer) {
	writer.WriteByte(126)

	stream := model.NewStream()
	stream.Write([]byte(s.Message))
	stream.WriteByte(10) // end message
	stream.WriteWordA(s.InterfaceId)

	buffer := stream.Flush()
	size := len(buffer)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(buffer)
}
