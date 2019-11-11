package outgoing

import "bufio"

type ConfigurationPacket struct {
	InterfaceId int
	State       int
}

func (c *ConfigurationPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(36)
	writer.WriteByte(byte(c.InterfaceId))
	writer.WriteByte(byte(c.InterfaceId >> 8))
	writer.WriteByte(byte(c.State))
}
