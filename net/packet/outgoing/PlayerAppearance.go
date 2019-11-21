package outgoing

import (
	"rsps/model"
)

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
	Target    model.PlayerInterface
	Equipment *model.ItemContainer
}

var EQUIPMENT_SLOTS = map[string]int{
	"head":   0,
	"cape":   1,
	"neck":   2,
	"weapon": 3,
	"chest":  4,
	"shield": 5,
	"legs":   7,
	"hands":  9,
	"feet":   10,
	"ring":   12,
	"ammo":   13,
}

func (p *PlayerAppearance) ToBytes() []byte {
	stream := model.NewStream()

	stream.WriteByte(0)   // gender
	stream.WriteByte(255) // prayer icon
	stream.WriteByte(255) // pk icon

	p.wordOrByte(stream, p.Equipment.Items[EQUIPMENT_SLOTS["head"]].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[EQUIPMENT_SLOTS["cape"]].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[EQUIPMENT_SLOTS["neck"]].ItemId)
	p.wordOrByte(stream, p.Equipment.Items[EQUIPMENT_SLOTS["weapon"]].ItemId)

	if p.Equipment.Items[EQUIPMENT_SLOTS["chest"]].ItemId > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[EQUIPMENT_SLOTS["chest"]].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[2])
	}

	p.wordOrByte(stream, p.Equipment.Items[EQUIPMENT_SLOTS["shield"]].ItemId)

	stream.WriteWord(0x100 + defaultAppearance[3]) //!isFullBody

	if p.Equipment.Items[EQUIPMENT_SLOTS["legs"]].ItemId > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[EQUIPMENT_SLOTS["legs"]].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[5])
	}

	stream.WriteWord(0x100 + defaultAppearance[1]) //isFullHelm

	if p.Equipment.Items[EQUIPMENT_SLOTS["hands"]].ItemId > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[EQUIPMENT_SLOTS["hands"]].ItemId))
	} else {
		stream.WriteWord(0x100 + defaultAppearance[4])
	}
	if p.Equipment.Items[EQUIPMENT_SLOTS["feet"]].ItemId > 1 {
		stream.WriteWord(0x200 + uint(p.Equipment.Items[EQUIPMENT_SLOTS["feet"]].ItemId))
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

	name := []byte(p.Target.GetName())
	stream.Write([]byte{0, 0, 0, 0, 0, 0, 0, name[0]}) //player name as int
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
