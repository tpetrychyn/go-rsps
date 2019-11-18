package incoming

import (
	"fmt"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
	"strconv"
	"strings"
)

type CommandsPacketHandler struct{}

func (c *CommandsPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	stream := model.NewStream()
	b := packet.ReadByte()
	for b != 10 {
		stream.WriteByte(b)
		b = packet.ReadByte()
	}

	parts := strings.Split(string(stream.Flush()), " ")
	command := parts[0]
	switch command {
	case "pos":
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("X: %d, Y: %d", player.Position.X, player.Position.Y)})
	case "item":
		if len(parts) == 1 { return }
		id, _ := strconv.Atoi(parts[1])
		amount := 1
		if len(parts) > 2 {
			amount, _ = strconv.Atoi(parts[2])
		}
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("adding item %d amount %d", id, amount)})
		player.Inventory.AddItem(id, amount)
	case "object":
		if len(parts) == 1 { return }
		id, _ := strconv.Atoi(parts[1])
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: id,
			Position: player.Position,
			Face:     0,
			Typ:      10,
			Player:   player,
		})
	}
}
