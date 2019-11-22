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
	writer.WriteByte(53)

	buffer := model.NewStream()
	buffer.WriteWord(s.InterfaceId)
	buffer.WriteWord(s.ItemContainer.Capacity)

	for _, v := range s.ItemContainer.Items {
		if v.Amount > 254 {
			buffer.WriteByte(255)
			buffer.WriteDWord_v2(v.Amount)
		} else {
			buffer.WriteByte(byte(v.Amount))
		}
		if v.ItemId > 0 {
			buffer.WriteWordBEA(uint(v.ItemId + 1))
		} else {
			buffer.WriteWordBEA(0)
		}
	}

	b := buffer.Flush()
	size := len(b)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(b)
}
