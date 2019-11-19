package entity

import (
	"fmt"
	"github.com/google/uuid"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"sync"
	"time"
)

type GroundItem struct {
	Position  *model.Position
	ItemId    int
	Amount    int
	Owner     *Player
	CreatedAt time.Time
}

type Region struct {
	Id uint16
	// TODO: convert all these maps to sync.Maps for thread safety
	Players           map[uuid.UUID]*Player
	GroundItems       *sync.Map
	WorldObjects      map[string]model.WorldObjectInterface // key is x-y as string
	MarkedForDeletion bool
}

func (r *Region) GetPlayersAsInterface() []model.PlayerInterface {
	var players = make([]model.PlayerInterface, 0)
	if r.Players == nil { // TODO: has crashed a few times nil on r.Players..
		return players
	}
	for _, v := range r.Players {
		players = append(players, v)
	}
	return players
}

func CreateRegion(id uint16) *Region {
	groundItems := new(sync.Map)
	if id == 12850 {
		groundItems.Store(uuid.New(), &GroundItem{
			Position: &model.Position{
				X: 3200,
				Y: 3200,
				Z: 0,
			},
			ItemId:    1351,
			Amount:    1,
			Owner:     nil,
			CreatedAt: time.Now(),
		})
	}

	return &Region{
		Id:           id,
		GroundItems:  groundItems,
		Players:      make(map[uuid.UUID]*Player),
		WorldObjects: make(map[string]model.WorldObjectInterface),
	}
}

func (r *Region) Tick() {
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if time.Now().Sub(g.CreatedAt) > 10*time.Second {
			// TODO: this works, bronze axe is simply spawned on lumby region creation atm
			//  adding items in constructor is good for global spawns
			r.RemoveGroundItemIdAtPosition(g.ItemId, g.Position)
			//delete(r.GroundItems, k)
			r.GroundItems.Delete(key)
		}

		if g.Owner != nil && time.Now().Sub(g.CreatedAt) > 5*time.Second {
			for _, player := range r.Players {
				if player.Id == g.Owner.Id {
					continue
				}
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
					Position:   g.Position,
					Player:     player,
					ItemId:     g.ItemId,
					ItemAmount: g.Amount,
				})
			}
			g.Owner = nil
		}

		return true
	})

	for _, obj := range r.WorldObjects {
		obj.Tick()
		if obj.ShouldRefresh() {
			for _, player := range r.Players {
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
					ObjectId: obj.GetObjectId(),
					Position: obj.GetPosition(),
					Face:     0,
					Typ:      10,
					Player:   player,
				})
			}
		}
	}

	var giLength = 0
	r.GroundItems.Range(func(key, value interface{}) bool {
		giLength++
		return true
	})
	if len(r.Players) == 0 && giLength == 0 {
		r.MarkedForDeletion = true
	}
}

func (r *Region) OnEnter(player *Player) {
	// We must add the player to all 9 regions entered on change
	// so that they get updates about regions around themselves
	r.Players[player.Id] = player
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if g.Owner != nil && g.Owner != player {
			return true
		} // only show it if nobody owns it or you own it
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
			Position:   g.Position,
			Player:     player,
			ItemId:     g.ItemId,
			ItemAmount: g.Amount,
		})
		return true
	})

	for _, obj := range r.WorldObjects {
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: obj.GetObjectId(),
			Position: obj.GetPosition(),
			Face:     0,
			Typ:      10,
			Player:   player,
		})
	}
}

func (r *Region) OnLeave(player *Player) {
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
			Position: g.Position,
			Player:   player,
			ItemId:   g.ItemId,
		})
		return true
	})
	delete(r.Players, player.Id)
}

func (r *Region) GetAdjacentIds() []uint16 {
	top := r.Id - 1
	bottom := r.Id + 1
	// includes r itself
	return []uint16{top - 256, top, top + 256, r.Id - 256, r.Id, r.Id + 256, bottom - 256, bottom, bottom + 256}
}

func (r *Region) FindGroundItemByPosition(id int, p *model.Position) *GroundItem {
	var groundItem *GroundItem
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if g.ItemId == id && g.Position.X == p.X && g.Position.Y == p.Y {
			groundItem = g
			return false
		}
		return true
	})
	return groundItem
}

func (r *Region) CreateGroundItemAtPosition(dropper *Player, item *model.Item, p *model.Position) {
	r.GroundItems.Store(uuid.New(), &GroundItem{
		Position:  p,
		ItemId:    item.ItemId,
		Amount:    item.Amount,
		Owner:     dropper,
		CreatedAt: time.Now(),
	})

	dropper.OutgoingQueue = append(dropper.OutgoingQueue, &outgoing.CreateGroundItemPacket{
		Position:   dropper.Position,
		Player:     dropper,
		ItemId:     item.ItemId,
		ItemAmount: item.Amount,
	})
}

func (r *Region) RemoveGroundItemIdAtPosition(id int, position *model.Position) {
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if g.ItemId == id && g.Position.X == position.X && g.Position.Y == position.Y {
			r.GroundItems.Delete(key)
			for _, p := range r.Players {
				p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
					Position: g.Position,
					Player:   p,
					ItemId:   g.ItemId,
				})
			}
			return false
		}
		return true
	})
}

func (r *Region) SetWorldObject(worldObject model.WorldObjectInterface) {
	key := fmt.Sprintf("%d-%d", worldObject.GetPosition().X, worldObject.GetPosition().Y)
	r.WorldObjects[key] = worldObject
	for _, p := range r.Players {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: worldObject.GetObjectId(),
			Position: worldObject.GetPosition(),
			Face:     0,
			Typ:      10,
			Player:   p,
		})
	}
}

func (r *Region) GetWorldObject(position *model.Position) model.WorldObjectInterface {
	key := fmt.Sprintf("%d-%d", position.X, position.Y)
	return r.WorldObjects[key]
}

func GetRegionIdByPosition(p *model.Position) uint16 {
	regionX := p.X >> 3
	regionY := p.Y >> 3
	return (regionX / 8 << 8) + (regionY / 8)
}
