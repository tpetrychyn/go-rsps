package packet

import (
	"bytes"
)

type Packet struct {
	Opcode  byte
	Size    uint16
	Buffer *bytes.Buffer
}

func (p *Packet) ReadBEShortA() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(lo) & 0xFF) << 8) | ((int(hi) - 128) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}

func (p *Packet) ReadBEShort() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(lo) & 0xFF) << 8) | (int(hi) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}