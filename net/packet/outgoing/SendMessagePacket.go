package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendMessagePacket struct {
	Message string
}

func (s *SendMessagePacket) Write(writer *bufio.Writer) {
	buffer := model.NewStream()
	buffer.Write([]byte(s.Message))
	buffer.WriteByte(10)

	payload := buffer.Flush()
	size := len(payload)
	writer.WriteByte(253)
	writer.WriteByte(byte(size))
	writer.Write(payload)
}
