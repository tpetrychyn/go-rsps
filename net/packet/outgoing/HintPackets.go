package outgoing

import (
	"bufio"
	"rsps/model"
)

type MobHintPacket struct {
	Type int
	Id   uint
}

func (f *MobHintPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(254)
	buffer := model.NewStream()
	buffer.WriteByte(byte(f.Type))
	buffer.WriteWord(f.Id)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)

	writer.Write(buffer.Flush())
}

type TileHintPacket struct {
	TilePosition int// (middle = 2; west = 3; east = 4; south = 5; north = 6)
	Position *model.Position
}

func (f *TileHintPacket) Write(writer *bufio.Writer) {
	if f.TilePosition < 2 || f.TilePosition > 6 {
		f.TilePosition = 2
	}
	writer.WriteByte(254)
	buffer := model.NewStream()
	buffer.WriteByte(byte(f.TilePosition))
	buffer.WriteWord(uint(f.Position.X))
	buffer.WriteWord(uint(f.Position.Y))
	buffer.WriteByte(0)
	writer.Write(buffer.Flush())
}

type ClearHintPacket struct {
}

func (f *ClearHintPacket) Write(writer *bufio.Writer) {
	writer.WriteByte(254)
	buffer := model.NewStream()
	buffer.WriteByte(0)
	buffer.WriteWord(255)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	writer.Write(buffer.Flush())
}

