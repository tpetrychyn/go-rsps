package packet

import "bufio"

type Packet struct {
	Opcode  byte
	Size    uint16
	Payload []byte
}

func (p *Packet) Write(writer *bufio.Writer) {
	writer.WriteByte(p.Opcode)
	writer.WriteByte(byte(p.Size >> 8))
	writer.WriteByte(byte(p.Size))
	writer.Write(p.Payload)
}

func (p *Packet) WriteAndFlush(writer *bufio.Writer) {
	writer.WriteByte(p.Opcode)
	writer.WriteByte(byte(p.Size >> 8))
	writer.WriteByte(byte(p.Size))
	writer.Write(p.Payload)

	writer.Flush()
}