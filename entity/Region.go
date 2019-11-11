package entity

import (
	"github.com/google/uuid"
	"math/rand"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"time"
)

type GroundItem struct {
	Position   *model.Position
	ItemId     int
	ItemAmount int
	Owner      *Player
	CreatedAt  time.Time
}

func (g *GroundItem) Tick() {
}

type Region struct {
	Id                uint16
	Players           map[uuid.UUID]*Player
	GroundItems       map[uuid.UUID]*GroundItem
	RemovedItems      []*GroundItem
	MarkedForDeletion bool
}

func CreateRegion(id uint16) *Region {
	groundItems := make(map[uuid.UUID]*GroundItem, 0)
	groundItems[uuid.New()] = &GroundItem{
		Position: &model.Position{
			X: 3200,
			Y: 3200,
			Z: 0,
		},
		ItemId:     rand.Intn(4500-1000) + 1000,
		ItemAmount: 10000,
		Owner:      nil,
		CreatedAt:  time.Now(),
	}
	return &Region{
		Id:           id,
		GroundItems:  groundItems,
		RemovedItems: make([]*GroundItem, 0),
		Players:      make(map[uuid.UUID]*Player),
	}
}

func (r *Region) Tick() {
	//for k, v := range r.GroundItems {
	//	if time.Now().Sub(v.CreatedAt) > 5*time.Second {
	//		r.RemovedItems = append(r.RemovedItems, v)
	//		delete(r.GroundItems, k)
	//	}
	//}

	for _, p := range r.Players {
		for _, i := range r.RemovedItems {
			p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
				Position: i.Position,
				Player:   p,
				ItemId:   i.ItemId,
			})
		}
	}

	r.RemovedItems = make([]*GroundItem, 0)
	if len(r.Players) == 0 && len(r.GroundItems) == 0 {
		r.MarkedForDeletion = true
	}
}

func (r *Region) OnEnter(player *Player) {
	r.Players[player.Id] = player
	for _, v := range r.GroundItems {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
			Position:   v.Position,
			Player:     player,
			ItemId:     v.ItemId,
			ItemAmount: v.ItemAmount,
		})
	}
}

func (r *Region) OnLeave(player *Player) {
	for _, v := range r.GroundItems {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
			Position: v.Position,
			Player:   player,
			ItemId:   v.ItemId,
		})
	}
	delete(r.Players, player.Id)
}

func (r *Region) GetAdjacentIds() []uint16 {
	//for _, v := range [][]int{{-1,-1,},{-1,0,},{-1,1},{0,-1},{0,1},{1,-1},{1,0},{1,1}} {
	//	regionId := GetRegionIdByPosition(&model.Position{
	//		X: uint16(int(player.Position.X) + (v[0]*64)),
	//		Y: uint16(int(player.Position.Y) + (v[1]*64)),
	//	})
	//	adjacent = append(adjacent, regionId)
	//	log.Printf("adjacent %d", regionId)
	//}
	top := r.Id - 1
	bottom := r.Id + 1
	return []uint16{top-256, top, top+256,r.Id-256,r.Id+256,bottom-256,bottom,bottom+256}
}

func GetRegionIdByPosition(p *model.Position) uint16 {
	regionX := p.X >> 3
	regionY := p.Y >> 3
	return (regionX / 8 << 8) + (regionY / 8)
}

