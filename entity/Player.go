package entity

import (
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/repository"
)

type OngoingAction interface {
	Tick()
}

type Player struct {
	*model.Movement
	MovementQueue *MovementQueue

	Id   int
	Name string

	Inventory          *Inventory
	Equipment          *Equipment
	Bank               *Bank
	SkillHelper        *SkillHelper
	OngoingAction      OngoingAction
	GlobalTickCount    int
	Region             *Region
	OutgoingQueue      []DownstreamMessage
	LoadedPlayers      []model.PlayerInterface
	LoadedNpcs         []model.NpcInterface
	UpdateFlag         *model.UpdateFlag
	DelayedPacket      func()
	DelayedDestination *model.Position
	LogoutRequested    bool
}

var SIDEBARS = []int{2423, 3917, 638, 3213, 1644, 5608, 1151,
	18128, 5065, 5715, 2449, 904, 147, 962}

func NewPlayer() *Player {
	spawn := &model.Position{X: 3200, Y: 3200}
	player := &Player{
		UpdateFlag: &model.UpdateFlag{},
		Movement: &model.Movement{
			Position:           spawn,
			LastKnownRegion:    spawn,
			PrimaryDirection:   model.None,
			SecondaryDirection: model.None,
			IsRunning:          true,
		},
		LoadedPlayers: make([]model.PlayerInterface, 0),
		LoadedNpcs:    make([]model.NpcInterface, 0),
	}
	player.MovementQueue = NewMovementQueue(player)
	player.Inventory = NewInventory(player)
	player.Equipment = NewEquipment(player)
	player.Bank = NewBank(player)
	player.SkillHelper = NewSkillHelper(player)

	return player
}

func (p *Player) SetId(id int) {
	p.Id = id
}

func (p *Player) LoadPlayer(name string) error {
	p.Name = name
	loadedItems, err := repository.InventoryRepositorySingleton.Load(name)
	if err == nil {
		p.Inventory.ItemContainer.Items = loadedItems
	}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Inventory.ItemContainer,
		InterfaceId:   model.INVENTORY_INTERFACE_ID,
	})

	loadedEquipment, err := repository.EquipmentRepositorySingleton.Load(name)
	if err == nil {
		p.Equipment.ItemContainer.Items = loadedEquipment
	}
	p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendItemContainerPacket{
		ItemContainer: p.Equipment.ItemContainer,
		InterfaceId:   model.EQUIPMENT_INTERFACE_ID,
	})

	loadedBank, err := repository.BankRepositorySingleton.Load(name)
	if err == nil {
		p.Bank.ItemContainer.Items = loadedBank
	}

	loadedSkills, err := repository.SkillRepositorySingleton.Load(name)
	if err == nil {
		p.SkillHelper.Skills = loadedSkills
	}
	for k, v := range p.SkillHelper.Skills {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SetSkillLevelPacket{
			SkillNum:   k,
			Level:      v.Level,
			Experience: v.Experience,
		})
	}

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

	if p.GlobalTickCount > 0 {
		p.GlobalTickCount--
	}

	if p.OngoingAction != nil {
		p.OngoingAction.Tick()
	}
}

func (p *Player) PostUpdate() {
	p.PrimaryDirection = model.None
	p.SecondaryDirection = model.None
	p.LastDirection = model.None
	p.UpdateFlag.Clear()
}

func (p *Player) Teleport(position *model.Position) {
	p.LastDirection = model.None
	p.PrimaryDirection = model.None
	p.Position = position
	p.LastKnownRegion = position
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

func (p *Player) GetId() int {
	return p.Id
}

func (p *Player) GetNearbyPlayers() []model.PlayerInterface {
	var players []model.PlayerInterface
	adjacentRegions := p.Region.GetAdjacentIds()
	for _, v := range adjacentRegions {
		r := world.GetRegion(v)
		for _, player := range r.GetPlayersAsInterface() {
			if p.GetId() == player.GetId() {
				continue
			}
			var found bool
			for _, addedPlayer := range players {
				if player.GetId() == addedPlayer.GetId() {
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

func (p *Player) GetLoadedPlayers() []model.PlayerInterface {
	return p.LoadedPlayers
}

func (p *Player) AddLoadedPlayer(pi model.PlayerInterface) {
	p.LoadedPlayers = append(p.LoadedPlayers, pi)
}

func (p *Player) RemoveLoadedPlayer(i int) {
	p.LoadedPlayers = append(p.LoadedPlayers[:i], p.LoadedPlayers[i+1:]...)
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetNearbyNpcs() []model.NpcInterface {
	var npcs []model.NpcInterface
	adjacentRegions := p.Region.GetAdjacentIds()
	for _, v := range adjacentRegions {
		r := world.GetRegion(v)
		for _, npc := range r.GetNpcsAsInterface() {
			var found bool
			for _, addedNpc := range npcs {
				if npc.GetId() == addedNpc.GetId() {
					found = true
					break
				}
			}
			if !found {
				npcs = append(npcs, npc)
			}
		}
	}
	return npcs
}
func (p *Player) GetLoadedNpcs() []model.NpcInterface {
	return p.LoadedNpcs
}
func (p *Player) AddLoadedNpc(n model.NpcInterface) {
	p.LoadedNpcs = append(p.LoadedNpcs, n)
}
func (p *Player) RemoveLoadedNpc(id int) {
	for k, v := range p.LoadedNpcs {
		if v.GetId() == id {
			p.LoadedNpcs = append(p.LoadedNpcs[:k], p.LoadedNpcs[k+1:]...)
			return
		}
	}
}

func (p *Player) GetInteractingWith() model.Character {
	return p.UpdateFlag.InteractingWith
}

func (p *Player) GetCurrentHitpoints() int {
	return p.SkillHelper.Skills[model.Hitpoints].Level
}

func (p *Player) GetMaxHitpoints() int {
	return p.SkillHelper.Skills[model.Hitpoints].GetLevelForExperience()
}

func (p *Player) TakeDamage(damage int) {
	p.SkillHelper.Skills[model.Hitpoints].Level -= damage
}

func (p *Player) GetMarkedForDeletion() bool {
	return p.LogoutRequested
}
//
//func (p *Player) String() string {
//	return "'"
//}
//
//func (p *Player) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
//	return nil, nil
//}
//func (p *Player) IsFalsy() bool {
//	return false
//}
//func (p *Player) Equals(x objects.Object) bool {
//	return false
//}
//
//func (p *Player) Copy() objects.Object {
//	return nil
//}
//func (p *Player) TypeName() string {
//	return ""
//}
//
//func (p *Player) IndexGet(index objects.Object) (objects.Object, error) {
//	val, _ := reflections.GetField(p, strings.Trim(index.String(), "\""))
//	obj, _ := objects.FromInterface(val)
//	return obj, nil
//}