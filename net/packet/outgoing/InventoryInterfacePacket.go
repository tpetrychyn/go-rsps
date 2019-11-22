package outgoing

import (
	"bufio"
	"rsps/model"
)

type InventoryInterfacePacket struct {
	InterfaceId        int
	SidebarInterfaceId int
}

func NewInventoryInterfacePacket(interfaceId, sidebarInterfaceId int) *InventoryInterfacePacket {
	return &InventoryInterfacePacket{
		InterfaceId:        interfaceId,
		SidebarInterfaceId: sidebarInterfaceId,
	}
}

func (i *InventoryInterfacePacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(248)
	stream.WriteWordA(uint(i.InterfaceId))
	stream.WriteWord(uint(i.SidebarInterfaceId))
	writer.Write(stream.Flush())
}
