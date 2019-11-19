package entity

import (
	"rsps/model"
	"sync"
)

var world *World

func WorldProvider() *World {
	if world == nil {
		world = CreateWorld()
	}
	return world
}

type World struct {
	Regions *sync.Map
}

func CreateWorld() *World {
	return &World{
		Regions: new(sync.Map),
	}
}

func (w *World) Tick() {
	//for {
	//<-time.After(600 * time.Millisecond)
	w.Regions.Range(func(key, value interface{}) bool {
		region := value.(*Region)
		if region.MarkedForDeletion {
			w.Regions.Delete(key)
			return true
		}
		region.Tick()
		return true
	})
	//}
}
func (w *World) GetRegion(id uint16) *Region {
	region, _ := w.Regions.LoadOrStore(id, CreateRegion(id))
	return region.(*Region)
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

func (w *World) SetWorldObject(worldObject model.WorldObjectInterface) {
	regionId := GetRegionIdByPosition(worldObject.GetPosition())
	w.GetRegion(regionId).SetWorldObject(worldObject)
}

func (w *World) GetWorldObject(position *model.Position) model.WorldObjectInterface {
	regionId := GetRegionIdByPosition(position)
	return w.GetRegion(regionId).GetWorldObject(position)
}
