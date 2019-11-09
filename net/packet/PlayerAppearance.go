package packet

import "rsps/model"

// TODO: Full helms and capes and such

var defaultAppearance = []uint{
	0, // gender
	0, // head
	18, // Torso
	26, // arms
	33, // hands
	36, // legs
	42, // feet
	10, // beard
	0, // hair colour
	0, // torso colour
	0, // legs colour
	0, // feet colour
	0, // skin colour
}

type PlayerAppearance struct {
	Hat uint
	Cape uint
	Amulet uint
	Weapon uint
	Chest uint
	Shield uint
	Legs uint
	Hands uint
	Feet uint
}

func (p *PlayerAppearance) ToBytes() []byte {
	stream := model.NewStream()

	stream.WriteByte(0)    // gender
	stream.WriteByte(3)    // prayer icon
	stream.WriteByte(0xFF) // pk icon

	p.wordOrByte(stream, p.Hat)
	p.wordOrByte(stream, p.Cape)
	p.wordOrByte(stream, p.Amulet)
	p.wordOrByte(stream, p.Weapon)

	if p.Chest > 1 {
		stream.WriteWord(0x200 + p.Chest)
	} else {
		stream.WriteWord(0x100 + defaultAppearance[2])
	}

	p.wordOrByte(stream, p.Shield)

	stream.WriteWord(0x100 + defaultAppearance[3]) //!isFullBody

	if p.Legs > 1 {
		stream.WriteWord(0x200 + p.Legs)
	} else {
		stream.WriteWord(0x100 + defaultAppearance[5])
	}

	stream.WriteWord(0x100 + defaultAppearance[1]) //isFullHelm

	if p.Hands > 1 {
		stream.WriteWord(0x200 + p.Hands)
	} else {
		stream.WriteWord(0x100 + defaultAppearance[4])
	}
	if p.Feet > 1 {
		stream.WriteWord(0x200 + p.Feet)
	} else {
		stream.WriteWord(0x100 + defaultAppearance[6])
	}

	stream.WriteWord(0x100 + defaultAppearance[7]) //gender && notFullHelm

	stream.WriteByte(defaultAppearance[8])
	stream.WriteByte(defaultAppearance[9])
	stream.WriteByte(defaultAppearance[10])
	stream.WriteByte(defaultAppearance[11])
	stream.WriteByte(defaultAppearance[12])

	stream.Write([]byte{0x328 >> 8, 0x328 & 0xFF})
	stream.Write([]byte{0x337 >> 8, 0x337 & 0xFF})
	stream.Write([]byte{0x333 >> 8, 0x333 & 0xFF})
	stream.Write([]byte{0x334 >> 8, 0x334 & 0xFF})
	stream.Write([]byte{0x335 >> 8, 0x335 & 0xFF})
	stream.Write([]byte{0x336 >> 8, 0x336 & 0xFF})
	stream.Write([]byte{0x338 >> 8, 0x338 & 0xFF})

	stream.Write([]byte{0, 0, 0, 0, 79, 120, 111, 6}) //player name as int
	stream.WriteByte(3)                               // combat level
	stream.Write([]byte{0, 0})                        // player skill level

	buffer := stream.Flush()
	updateSize := len(buffer)-1
	out := append([]byte{byte(^updateSize)}, buffer...)
	return out
}

func (p *PlayerAppearance) wordOrByte(stream *model.Stream, slot uint) {
	if slot > 1 {
		stream.WriteWord(0x200 + slot)
	} else {
		stream.WriteByte(0)
	}
}