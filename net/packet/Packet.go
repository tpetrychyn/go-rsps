package packet

import (
	"bytes"
)

type Packet struct {
	Opcode  byte
	Size    uint16
	Payload []byte
	Buffer  *bytes.Buffer
}

func (p *Packet) ReadLEShortA() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(lo) & 0xFF) << 8) | ((int(hi) - 128) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}

func (p *Packet) ReadLEShort() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(lo) & 0xFF) << 8) | (int(hi) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}

func (p *Packet) ReadShortA() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(hi) & 0xFF) << 8) | (int(lo - 128) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}

func (p *Packet) ReadShort() uint16 {
	hi, _ := p.Buffer.ReadByte()
	lo, _ := p.Buffer.ReadByte()
	value := ((int(hi) & 0xFF) << 8) | (int(lo) & 0xFF)
	if int(uint16(value)) > 32767 {
		value -= 0x10000
	}
	return uint16(value)
}

func (p *Packet) ReadByte() byte {
	value, _ := p.Buffer.ReadByte()
	return value
}
