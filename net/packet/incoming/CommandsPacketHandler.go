package incoming

import (
	"context"
	"fmt"
	"github.com/mattn/anko/vm"
	"log"
	"reflect"
	"rsps/entity"
	"rsps/handler"
	"rsps/model"
	"rsps/net/packet"
	"rsps/net/packet/outgoing"
	"rsps/util"
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
		if len(parts) == 1 {
			return
		}
		var id int
		itemByName := util.GetItemDefinitionByName(parts[1], false)
		if itemByName == nil {
			id, _ = strconv.Atoi(parts[1])
		} else {
			id = itemByName.ID
		}
		amount := 1
		if len(parts) > 2 {
			amount, _ = strconv.Atoi(parts[2])
		}

		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("adding item %d amount %d", id, amount)})
		player.Inventory.AddItem(id, amount)

	case "object":
		if len(parts) == 1 {
			return
		}
		id, _ := strconv.Atoi(parts[1])
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: id,
			Position: player.Position,
			Face:     0,
			Typ:      10,
			Player:   player,
		})

	case "tele":
		if len(parts) == 2 {
			return
		}
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[1])
		player.Teleport(&model.Position{X: uint16(x), Y: uint16(y)})

	case "region":
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("You are in region: %v", entity.GetRegionIdByPosition(player.Position))})

	case "bank":
		player.Bank.OpenBank()

	case "r":
		handler.LoadScripts()
	default:
		f := handler.CommandObservers[command]
		if f != nil {
			if f, ok := f.(func(ctx context.Context, value reflect.Value, value2 reflect.Value) (reflect.Value, reflect.Value)); ok {
				e := handler.WorldModule()
				err := e.Define("exec", func() {
					_, err := f(context.Background(), reflect.ValueOf(player), reflect.ValueOf(parts))
					if !err.IsNil() {
						log.Printf("err %s", err)
					}
				})
				_, err = vm.Execute(e, nil, `exec()`)
				if err != nil {
					log.Printf("err %s", err.Error())
				}
			}
		}
	}

}
