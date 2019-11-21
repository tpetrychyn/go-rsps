package outgoing

import (
	"bufio"
	"rsps/model"
)

type PlayerDetailsPacket struct {
	Pid int
}

func (p *PlayerDetailsPacket) Write(writer *bufio.Writer) {
	stream := model.NewStream()
	stream.WriteByte(249)
	stream.WriteByte(1 + 128)
	stream.WriteWordBEA(uint(p.Pid))

	writer.Write(stream.Flush())
}
