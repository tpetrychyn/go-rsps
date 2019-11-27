package incoming

import (
	"context"
	"github.com/mattn/anko/vm"
	"log"
	"reflect"
	"rsps/entity"
	"rsps/handler"
	"rsps/model"
	"rsps/net/packet"
)

type ObjectActionPacket struct{}

func (e *ObjectActionPacket) HandlePacket(player *entity.Player, packet *packet.Packet) {
	switch packet.Opcode {
	case OBJECT_ACTION_ONE_OPCODE:
		e.handleObjectActionOne(player, packet)
	}
}

func (e *ObjectActionPacket) handleObjectActionOne(player *entity.Player, packet *packet.Packet) {
	objectX := packet.ReadLEShortA()
	objectId := packet.ReadShort()
	objectY := packet.ReadShortA()

	e.handleObjectActionOneInternal(player, objectX, objectY, objectId)
}

func (e *ObjectActionPacket) handleObjectActionOneInternal(player *entity.Player, x, y, id uint16) {
	if player.DelayedDestination != nil {
		player.DelayedPacket = func() {
			e.handleObjectActionOneInternal(player, x, y, id)
		}
		return
	}

	objPosition := &model.Position{X: x, Y: y}
	object := entity.WorldProvider().GetWorldObject(objPosition)
	if object == nil {
		object = &model.DefaultWorldObject{
			ObjectId: int(id),
			Position: objPosition,
			Face:     0,
			Type:     10,
		}
	}

	player.UpdateFlag.SetFacePosition(objPosition)

	f := handler.ObjectObservers[int(id)]
	if f != nil {
		if f, ok := f.(func(ctx context.Context, value reflect.Value, value2 reflect.Value) (reflect.Value, reflect.Value)); ok {
			e := handler.WorldModule()
			err := e.Define("exec", func() {
				_, err := f(context.Background(), reflect.ValueOf(player), reflect.ValueOf(object))
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
	//handler.LoadScripts(player)

	//if handler.ObjectObservers[int(id)] != nil {
	//	f := handler.ObjectObservers[int(id)].Function
	//	c := handler.ObjectObservers[int(id)].CompiledScript
	//
	//	o, err := objects.FromInterface(player)
	//	if err != nil {
	//		log.Printf("err %s", err.Error())
	//		return
	//	}
	//	c.Set("player", o)
	//	c.Set("object", o)
	//	err = c.Set("execute", f)
	//	if err != nil {
	//		log.Printf("err %s", err.Error())
	//		return
	//	}
	//	if err := c.Run(); err != nil {
	//		log.Printf("err %s", err.Error())
	//	}
	//	return
	//}

	switch id {
	case 2213:
		player.Bank.OpenBank()
		return
	}

	if handler.ObjectIsWoodcuttingTree(int(id)) {
		handler.StartWoodcutting(int(id), &model.Position{X: x, Y: y}, player)
		return
	}

	log.Printf("Object Click1: x %+v, y %+v, id %+v", x, y, id)
}
