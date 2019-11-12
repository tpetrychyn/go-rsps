package outgoing

import "bufio"

type LogoutPacket struct {}

func (l *LogoutPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(109)
}
