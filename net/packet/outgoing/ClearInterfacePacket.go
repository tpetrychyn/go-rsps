package outgoing

import "bufio"

type ClearInterfacePacket struct {}

func (c *ClearInterfacePacket) Write(writer *bufio.Writer) {
	writer.WriteByte(219)
}
