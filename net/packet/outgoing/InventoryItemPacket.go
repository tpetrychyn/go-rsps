package outgoing

import (
	"bufio"
	"rsps/model"
)

type InterfaceItemPacket struct {
	InterfaceId int
	Slot        int
	*model.Item
}

func (i *InterfaceItemPacket) Write(writer *bufio.Writer) {
	if i.Item == nil {
		return
	}
	writer.WriteByte(34)
	payload := model.NewStream()
	payload.WriteWord(uint(i.InterfaceId))
	payload.WriteByte(byte(i.Slot))
	if i.ItemId > 0 {
		payload.WriteWord(uint(i.ItemId + 1))
	} else {
		payload.WriteWord(0)
	}

	if i.Amount > 254 {
		payload.WriteByte(255)
		payload.WriteInt(i.Amount)
	} else {
		payload.WriteByte(byte(i.Amount))
	}
	b := payload.Flush()
	size := len(b)
	writer.WriteByte(byte(size >> 8))
	writer.WriteByte(byte(size))
	writer.Write(b)
}
