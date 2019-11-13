package outgoing

import "bufio"

type FlashSidebarPacket struct {
	SidebarId int
}

func (f *FlashSidebarPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(24)
	writer.WriteByte(byte(f.SidebarId + 128))
}
