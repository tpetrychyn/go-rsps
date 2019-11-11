package outgoing

import "bufio"

type SidebarInterfacePacket struct {
	MenuId int
	Form   int
}

func (s *SidebarInterfacePacket) Write(writer *bufio.Writer) {
	writer.WriteByte(71)
	writer.WriteByte(byte(s.Form >> 8))
	writer.WriteByte(byte(s.Form))
	writer.WriteByte(byte(s.MenuId + 128))
}
