package incoming

import (
	"fmt"
	"rsps/entity"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
	"strconv"
	"strings"
)

type ChatPacketHandler struct {}

func (e *ChatPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	_ = packet.ReadByteS()
	_ = packet.ReadByteS()
	size := packet.Size - 2
	word := make([]byte, size)
	for i:=int(size-1);i>=0;i-- {
		word[i] = packet.ReadByte() - 128
	}

	chat := e.unpack(word, int(size))
	if strings.HasPrefix(chat, "$") {
		if strings.Contains(chat, "item") {
			parts := strings.Split(chat, " ")
			id, _ := strconv.Atoi(parts[1])
			amount, _ := strconv.Atoi(parts[2])
			player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message:fmt.Sprintf("adding item %d amount %d", id, amount)})
			player.Inventory.AddItem(id, amount)
		}
	}
}

var xlateTable = []rune{ ' ', 'e', 't', 'a', 'o', 'i', 'h', 'n',
	's', 'r', 'd', 'l', 'u', 'm', 'w', 'c', 'y', 'f', 'g', 'p', 'b',
	'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2', '3', '4', '5', '6',
	'7', '8', '9', ' ', '!', '?', '.', ',', ':', ';', '(', ')', '-',
	'&', '*', '\\', '\'', '@', '#', '+', '=', '\243', '$', '%', '"',
	'[', ']' }

func (e *ChatPacketHandler) unpack(data []byte, size int) string {
	decodeBuf := make([]rune, 4096)
	idx := 0
	highNibble := -1
	for i:=0;i<size*2;i++ {
		val := (int(data[i/2]) >> uint(4 - 4 * (i % 2))) & 0xf
		if highNibble == -1 {
			if val < 13 {
				decodeBuf[idx] = xlateTable[val]
				idx += 1
			} else {
				highNibble = int(val)
			}
		} else {
			decodeBuf[idx] = xlateTable[(highNibble << 4) + int(val) - 195]
			idx += 1
			highNibble = -1
		}
	}
	return string(decodeBuf)
}
