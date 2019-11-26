package entity

import (
	"rsps/model"
	"sync"
)

var world *World

func WorldProvider() *World {
	if world == nil {
		world = CreateWorld()
		world.AddNpc(2, &model.Position{
			X: 3200,
			Y: 3200,
			Z: 0,
		})
	}
	return world
}

type World struct {
	Regions *sync.Map
	Players *sync.Map
	Npcs    *sync.Map
}

func CreateWorld() *World {
	return &World{
		Regions: new(sync.Map),
		Players: new(sync.Map),
		Npcs:    new(sync.Map),
	}
}

func (w *World) Tick() {
	w.Regions.Range(func(key, value interface{}) bool {
		region := value.(*Region)
		if region.MarkedForDeletion {
			w.Regions.Delete(key)
			return true
		}
		region.Tick()
		return true
	})
}

func (w *World) PostUpdate() {
	w.Regions.Range(func(key, value interface{}) bool {
		region := value.(*Region)
		region.PostUpdate()
		return true
	})
}

func (w *World) GetRegion(id uint16) *Region {
	region, _ := w.Regions.LoadOrStore(id, CreateRegion(id))
	return region.(*Region)
}

func (w *World) AddNpc(npcType int, position *model.Position) {
	region := w.GetRegion(GetRegionIdByPosition(position))
	for id := 0; id < 2000; id++ {
		_, ok := w.Npcs.Load(id)
		if !ok {
			npc := NewNpc(id)
			w.Npcs.Store(id, npc)
			region.Npcs.Store(id, npc)
			return
		}
	}
}

func (w *World) RemoveNpc(id int) {
	n, ok := w.Npcs.Load(id)
	if !ok {
		return
	}
	npc := n.(*Npc)
	region := w.GetRegion(GetRegionIdByPosition(npc.Position))
	region.Npcs.Delete(id)
	w.Npcs.Delete(id)
}

func (w *World) AddPlayer(player *Player) {
	for id := 1; id < 2000; id++ {
		_, ok := w.Players.Load(id)
		if !ok {
			player.Id = id
			player.Teleport(player.Position)
			w.Players.Store(id, player)
			return
		}
	}
}

func (w *World) RemovePlayer(id int) {
	p, ok := w.Players.Load(id)
	if !ok {
		return
	}
	player := p.(*Player)
	regionId := GetRegionIdByPosition(player.Position)
	region := w.GetRegion(regionId)

	adj := region.GetAdjacentIds()
	for _, x := range adj {
		world.GetRegion(x).OnLeave(player)
	}
	w.Players.Delete(id)
}

func (w *World) AddPlayerToRegion(player *Player) {
	previousRegion := player.Region
	regionId := GetRegionIdByPosition(player.Position)
	region := w.GetRegion(regionId)
	player.Region = region

	newAdj := region.GetAdjacentIds()
	var prevAdj []uint16
	if previousRegion != nil {
		prevAdj = previousRegion.GetAdjacentIds()
	}

	for _, x := range newAdj {
		var exists bool
		for _, y := range prevAdj {
			if x == y {
				exists = true
			}
		}
		if !exists {
			r := world.GetRegion(x)
			r.OnEnter(player)
		}
	}

	for _, x := range prevAdj {
		var exists bool
		for _, y := range newAdj {
			if x == y {
				exists = true
			}
		}
		if !exists {
			r := world.GetRegion(x)
			r.OnLeave(player)
		}
	}
}

func (w *World) GetRegionForPlayer(player *Player) *Region {
	regionId := GetRegionIdByPosition(player.Position)
	return w.GetRegion(regionId)
}

func (w *World) AddWorldObject(id int, position *model.Position) {
	obj := &model.DefaultWorldObject{ObjectId: id, Position: position}
	regionId := GetRegionIdByPosition(position)
	w.GetRegion(regionId).SetWorldObject(obj)
}

func (w *World) SetWorldObject(worldObject model.WorldObjectInterface) {
	regionId := GetRegionIdByPosition(worldObject.GetPosition())
	w.GetRegion(regionId).SetWorldObject(worldObject)
}

func (w *World) GetWorldObject(position *model.Position) model.WorldObjectInterface {
	regionId := GetRegionIdByPosition(position)
	return w.GetRegion(regionId).GetWorldObject(position)
}
