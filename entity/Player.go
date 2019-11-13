package entity

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/util"
	"time"
)

type Player struct {
	Id   uuid.UUID
	Name string
	*model.Movement
	MovementQueue   *MovementQueue
	Region          *Region
	Inventory       *model.ItemContainer
	Equipment       *model.ItemContainer
	OutgoingQueue   []DownstreamMessage
	LogoutRequested bool
}

var SIDEBARS = []int{2423, 3917, 638, 3213, 1644, 5608, 1151,
	18128, 5065, 5715, 2449, 904, 147, 962}

func NewPlayer() *Player {
	spawn := &model.Position{
		X: 3200,
		Y: 3200,
	}
	player := &Player{
		Id:        uuid.New(), // TODO: Load this from database or something
		Inventory: model.NewItemContainer(28),
		Equipment: model.NewItemContainer(14),
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			SecondaryDirection: model.None,
			IsRunning:          true,
		},
	}
	player.MovementQueue = NewMovementQueue(player)

	for k, v := range SIDEBARS {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SidebarInterfacePacket{
			MenuId: k,
			Form:   v,
		})
	}

	return player
}

type PlayerSave struct {
	Position *model.Position
	Inventory []*model.Item `json:"inventory"`
}

func (p *Player) SavePlayer() {
	save := &PlayerSave{Position: p.Position, Inventory: p.Inventory.Items}

	file, _ := json.MarshalIndent(save, "", " ")

	_ = ioutil.WriteFile(fmt.Sprintf("./players/%s.json", p.Name), file, 0644)
}

func (p *Player) LoadPlayer(name string) error {
	fname := fmt.Sprintf("./players/%s", name) + ".json"
	log.Printf("%+v", fname)
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	var playerSave PlayerSave
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &playerSave)
	if err != nil {
		return err
	}
	if len(playerSave.Inventory) > 0 {
		p.Inventory.Items = playerSave.Inventory
	}

	if playerSave.Position != nil {
		p.Position = playerSave.Position
		p.LastKnownRegion = playerSave.Position
	}

	p.Name = name

	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Inventory,
		InterfaceId:   3214,
	})

	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Welcome to TaylorScape"})
	log.Printf("loaded %s", fname)
	return nil
}

func (p *Player) PostUpdate() {
	p.PrimaryDirection = model.None
	p.SecondaryDirection = model.None
	p.LastDirection = model.None
}

var t = time.Now()

func (p *Player) Tick() {
	p.MovementQueue.Tick()
	p.SavePlayer()
}

func (p *Player) EquipItem(invSlot, itemId uint16) {
	invItem := p.Inventory.FindItem(int(itemId))
	if invItem == nil {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have that item."})
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.InventoryItemPacket{Slot: int(invSlot), Item: &model.Item{}})
		return
	}

	def := util.GetItemDefinition(int(itemId))
	slot := outgoing.EQUIPMENT_SLOTS[def.Equipment.Slot]

	p.Equipment.Items[slot] = invItem
	p.Inventory.Items[invSlot] = &model.Item{}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: int(invSlot),
		Item: &model.Item{},
	})

	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Equipment,
		InterfaceId:   model.EQUIPMENT_INTERFACE_ID,
	})
}

func (p *Player) Teleport(position *model.Position) {
	p.LastDirection = model.None
	p.PrimaryDirection = model.None
	p.Position = position
	p.LastKnownRegion = p.Position
	p.MovementQueue.Clear()
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.MapRegionPacket{Position: p.Position})
	p.OutgoingQueue = append(p.OutgoingQueue, outgoing.NewPlayerUpdatePacket(p).SetUpdateRequired(true).SetTyp(outgoing.Teleport))
}

func (p *Player) GetEquipment() *model.ItemContainer {
	return p.Equipment
}

func (p *Player) AddItem(id, amount int) {
	slot := p.Inventory.AddItem(id, amount)
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.InventoryItemPacket{
		Slot: slot,
		Item: &model.Item{
			ItemId: id,
			Amount: 1,
		},
	})
}
