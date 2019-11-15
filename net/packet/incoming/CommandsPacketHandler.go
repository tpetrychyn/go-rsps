package incoming

import (
	"fmt"
	"log"
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
		var amount int
		if len(parts) > 2 {
			amount, _ = strconv.Atoi(parts[2])
		}
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("adding item %d amount %d", id, amount)})
		player.Inventory.AddItem(id, amount)
	}
	log.Printf("command %+v", string(stream.Flush()))
}
