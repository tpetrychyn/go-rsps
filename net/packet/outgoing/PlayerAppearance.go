package outgoing

import "rsps/model"

// TODO: Full helms and capes and such

var defaultAppearance = []uint{
	0,  // gender
	0,  // head
	18, // Torso
	26, // arms
	33, // hands
	36, // legs
	42, // feet
	10, // beard
	0,  // hair colour
	0,  // torso colour
	0,  // legs colour
	0,  // feet colour
	0,  // skin colour
}

type PlayerAppearance struct {
	Equipment *model.ItemContainer
	Hat       int
	Cape      int
	Amulet    int
	Weapon    int
	Chest     int
	Shield    int
	Legs      int
	Hands     int
	Feet      int
}

const (
	HAT_SLOT    = 0
	CAPE_SLOT   = 1
	AMULET_SLOT = 2
	WEAPON_SLOT = 3
	CHEST_SLOT  = 4
	SHIELD_SLOT = 5
	LEGS_SLOT   = 7
	HANDS_SLOT  = 9
	FEET_SLOT   = 10
	RING_SLOT   = 12
	ARROW_SLOT  = 13
)

func (p *PlayerAppearance) ToBytes() []byte {
	stream := model.NewStream()

	stream.WriteByte(0)   // gender
	stream.WriteByte(255) // prayer icon
	stream.WriteByte(255) // pk icon

	p.wordOrByte(stream, p.Equipment.Items[HAT_SLOT].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[CAPE_SLOT].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[AMULET_SLOT].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[WEAPON_SLOT].ItemId)

	if p.Chest > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[CHEST_SLOT].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[2])
	}

	p.wordOrByte(stream, p.Equipment.Items[SHIELD_SLOT].ItemId)

	stream.WriteWord(0x100 + defaultAppearance[3]) //!isFullBody

	if p.Legs > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[LEGS_SLOT].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[5])
	}

	stream.WriteWord(0x100 + defaultAppearance[1]) //isFullHelm

	if p.Hands > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[HANDS_SLOT].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[4])
	}
	if p.Feet > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[FEET_SLOT].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[6])
	}

	stream.WriteWord(0x100 + defaultAppearance[7]) //gender && notFullHelm

	stream.WriteByte(byte(defaultAppearance[8]))
	stream.WriteByte(byte(defaultAppearance[9]))
	stream.WriteByte(byte(defaultAppearance[10]))
	stream.WriteByte(byte(defaultAppearance[11]))
	stream.WriteByte(byte(defaultAppearance[12]))

	stream.Write([]byte{0x328 >> 8, 0x328 & 0xFF})
	stream.Write([]byte{0x337 >> 8, 0x337 & 0xFF})
	stream.Write([]byte{0x333 >> 8, 0x333 & 0xFF})
	stream.Write([]byte{0x334 >> 8, 0x334 & 0xFF})
	stream.Write([]byte{0x335 >> 8, 0x335 & 0xFF})
	stream.Write([]byte{0x336 >> 8, 0x336 & 0xFF})
	stream.Write([]byte{0x338 >> 8, 0x338 & 0xFF})

	stream.Write([]byte{0, 0, 1, 168, 251, 9, 73, 127}) //player name as int
	stream.WriteByte(3)                                 // combat level
	stream.Write([]byte{0, 0})                          // player skill level

	buffer := stream.Flush()
	updateSize := len(buffer) - 1
	out := append([]byte{byte(^updateSize)}, buffer...)
	return out
}

func (p *PlayerAppearance) wordOrByte(stream *model.Stream, slot int) {
	if slot > 1 {
		stream.WriteWord(0x200 + uint(slot))
	} else {
		stream.WriteByte(0)
	}
}
