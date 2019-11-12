package outgoing

import (
	"bufio"
	"rsps/model"
)

type SendItemContainerPacket struct {
	ItemContainer *model.ItemContainer
	InterfaceId   uint
}

func (s *SendItemContainerPacket) Write(writer *bufio.Writer) {
	// TODO: This does not work
	writer.WriteByte(53)
	buffer := model.NewStream()

	buffer.WriteWord(s.InterfaceId)
	buffer.WriteWord(s.ItemContainer.Capacity)

	for _, _ = range s.ItemContainer.Items {
		buffer.WriteByte(0)
		buffer.WriteByte(0)
		buffer.WriteByte(0)
	}

	b := buffer.Flush()
	size := len(b)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(b)
}
