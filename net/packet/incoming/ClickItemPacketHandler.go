package incoming

import (
	"context"
	"github.com/mattn/anko/vm"
	"log"
	"reflect"
	"rsps/entity"
	"rsps/handler"
	"rsps/net/packet"
)

type ClickItemPacketHandler struct {}

func (c *ClickItemPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	_ = packet.ReadLEShortA() //interfaceId
	slot := packet.ReadShortA()
	id := packet.ReadLEShort()

	//log.Printf("clicked slot %v id %v", slot, id)

	handlers := handler.ItemClickObservers[int(id)]
	if handlers != nil {
		for _, f := range handlers {
			if f, ok := f.(func(ctx context.Context, player reflect.Value, id reflect.Value, slot reflect.Value) (reflect.Value, reflect.Value)); ok {
				e := handler.WorldModule()
				err := e.Define("exec", func() {
					_, err := f(context.Background(), reflect.ValueOf(player), reflect.ValueOf(id), reflect.ValueOf(slot))
					if !err.IsNil() {
						log.Printf("err %s", err)
					}
				})
				_, err = vm.Execute(e, nil, `exec()`)
				if err != nil {
					log.Printf("err %s", err.Error())
				}
			} else {
				log.Printf("bound itemClicks should have params (player, id, slot)")
			}
		}
	}
}
