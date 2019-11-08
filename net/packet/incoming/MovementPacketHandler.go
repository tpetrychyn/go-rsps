package incoming

import (
	"bytes"
	"encoding/binary"
	"log"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
)

type Point struct {
	X int8
	Y int8
}

type MovementPacketHandler struct {}

func (m *MovementPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	buffer := bytes.NewBuffer(packet.Payload)
	steps := int((packet.Size - 5) / 2)
	path := make([]Point, steps)
	var firstStepX uint16
	_ = binary.Read(buffer, binary.LittleEndian, &firstStepX)
	for i := 0; i < steps; i++ {
		var point Point
		_ = binary.Read(buffer, binary.BigEndian, &point)
		path[i] = point
	}
	var firstStepY uint16
	_ = binary.Read(buffer, binary.LittleEndian, &firstStepY)
	log.Printf("firstStep %+v %+v", firstStepX+128, firstStepY)
	log.Printf("Path %+v", path)

	positions := make([]*model.Position, steps+1)
	positions[0] = &model.Position{
		X: firstStepX,
		Y: firstStepY,
	}

	for i := 0; i < steps; i++ {
		positions[i+1] = &model.Position{
			X: firstStepX + uint16(path[i].X),
			Y: firstStepY + uint16(path[i].Y),
			Z: 0,
		}
	}

	for _, v := range positions {
		log.Printf("position %+v", v)
	}
}
