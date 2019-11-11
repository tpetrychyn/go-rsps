package entity

import (
	"time"
)

var world *World
func WorldProvider() *World {
	if world == nil {
		world = CreateWorld()
	}
	return world
}

type World struct {
	Regions map[uint16]*Region
}

func CreateWorld() *World {
	return &World{
		Regions: make(map[uint16]*Region),
	}
}

func (w *World) Tick() {
	for {
		<- time.After(100 * time.Millisecond)
		for k, v := range w.Regions {
			v.Tick()
			if v.MarkedForDeletion {
				delete(w.Regions, k)
			}
		}
	}
}

func (w *World) GetRegion(id uint16) *Region {
	if w.Regions[id] == nil {
		w.Regions[id] = CreateRegion(id)
	}
	return w.Regions[id]
}

func (w *World) AddPlayerToRegion(player *Player) {
	previousRegion := player.Region
	regionId := GetRegionIdByPosition(player.Position)
	region := w.GetRegion(regionId)
	if previousRegion != nil {
		previousRegion.OnLeave(player)
	}

	player.Region = region
	region.OnEnter(player)
}

func (w *World) GetRegionForPlayer(player *Player) *Region {
	regionId := GetRegionIdByPosition(player.Position)
	return w.Regions[regionId]
}

