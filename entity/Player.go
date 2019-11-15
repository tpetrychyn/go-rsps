package entity

import (
	"github.com/google/uuid"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type Player struct {
	*model.Movement

	Id              uuid.UUID
	Name            string
	MovementQueue   *MovementQueue
	Region          *Region
	Inventory       *Inventory
	Equipment       *Equipment
	OutgoingQueue   []DownstreamMessage
	LoadedPlayers   map[uuid.UUID]model.PlayerInterface
	UpdateFlag      *model.UpdateFlag
	DelayedPacket   func()
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
		Id:         uuid.New(), // TODO: Load this from database or something
		UpdateFlag: &model.UpdateFlag{},
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			SecondaryDirection: model.None,
			IsRunning:          true,
		},
		LoadedPlayers: make(map[uuid.UUID]model.PlayerInterface),
	}
	player.MovementQueue = NewMovementQueue(player)
	player.Inventory = NewInventory(player)
	player.Equipment = NewEquipment(player)

	return player
}

func (p *Player) LoadPlayer(name string) error {
	p.Name = name
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Inventory.ItemContainer,
		InterfaceId:   model.INVENTORY_INTERFACE_ID,
	})
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Equipment.ItemContainer,
		InterfaceId:   model.EQUIPMENT_INTERFACE_ID,
	})
	for k, v := range SIDEBARS {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SidebarInterfacePacket{
			MenuId: k,
			Form:   v,
		})
	}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Welcome to TaylorScape"})
	return nil
}

func (p *Player) Tick() {
	p.MovementQueue.Tick()
}

func (p *Player) PostUpdate() {
	p.PrimaryDirection = model.None
	p.SecondaryDirection = model.None
	p.LastDirection = model.None
}

func (p *Player) Teleport(position *model.Position) {
	p.LastDirection = model.None
	p.PrimaryDirection = model.None
	p.Position = position
	p.LastKnownRegion = p.Position
	p.MovementQueue.Clear()
	p.UpdateFlag.NeedsPlacement = true
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.MapRegionPacket{Position: p.Position})
	world.AddPlayerToRegion(p)
}

func (p *Player) GetEquipmentItemContainer() *model.ItemContainer {
	return p.Equipment.ItemContainer
}

func (p *Player) GetUpdateFlag() *model.UpdateFlag {
	return p.UpdateFlag
}

func (p *Player) GetPlayerId() uuid.UUID {
	return p.Id
}

func (p *Player) GetNearbyPlayers() []model.PlayerInterface {
	var players []model.PlayerInterface
	adjacentRegions := p.Region.GetAdjacentIds()
	for _, v := range adjacentRegions {
		r := world.Regions[v]
		if r == nil {
			continue
		}
		for _, player := range r.GetPlayersAsInterface() {
			if p.GetPlayerId() == player.GetPlayerId() {
				continue
			}
			var found bool
			for _, addedPlayer := range players {
				if player.GetPlayerId() == addedPlayer.GetPlayerId() {
					found = true
					break
				}
			}
			if !found {
				players = append(players, player)
			}
		}
	}
	return players
}

func (p *Player) GetLoadedPlayers() map[uuid.UUID]model.PlayerInterface {
	return p.LoadedPlayers
}

func (p *Player) AddLoadedPlayer(pi model.PlayerInterface) {
	p.LoadedPlayers[pi.GetPlayerId()] = pi
}

func (p *Player) RemoveLoadedPlayer(id uuid.UUID) {
	delete(p.LoadedPlayers, id)
}

func (p *Player) GetName() string {
	return p.Name
}
