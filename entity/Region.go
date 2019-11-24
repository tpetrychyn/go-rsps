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
	Id                uint16
	Players           *sync.Map
	Npcs              *sync.Map
	GroundItems       *sync.Map
	WorldObjects      *sync.Map
	MarkedForDeletion bool
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
		Players:      new(sync.Map),
		Npcs:         new(sync.Map),
		WorldObjects: new(sync.Map),
	}
}

func (r *Region) PostUpdate() {
	r.Npcs.Range(func(key, value interface{}) bool {
		npc := value.(*Npc)
		npc.PostUpdate()
		return true
	})
}

func (r *Region) Tick() {
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if time.Now().Sub(g.CreatedAt) > 2*time.Minute {
			r.RemoveGroundItemIdAtPosition(g.ItemId, g.Position)
			r.GroundItems.Delete(key)
		}

		if g.Owner != nil && time.Now().Sub(g.CreatedAt) > 1*time.Minute {
			r.Players.Range(func(key, value interface{}) bool {
				player := value.(*Player)
				if player.Name == g.Owner.Name {
					return true
				}
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
					Position:   g.Position,
					Player:     player,
					ItemId:     g.ItemId,
					ItemAmount: g.Amount,
				})
				return true
			})
			g.Owner = nil
		}
		return true
	})

	r.WorldObjects.Range(func(key, value interface{}) bool {
		obj := value.(model.WorldObjectInterface)
		obj.Tick()
		if obj.ShouldRefresh() {
			r.Players.Range(func(key, value interface{}) bool {
				player := value.(*Player)
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
					ObjectId: obj.GetObjectId(),
					Position: obj.GetPosition(),
					Face:     0,
					Typ:      10,
					Player:   player,
				})
				return true
			})
		}
		return true
	})

	r.Npcs.Range(func(key, value interface{}) bool {
		npc := value.(*Npc)
		npc.Tick()
		return true
	})

	// TODO: Hate having to loop just to get count...
	var giLength = 0
	r.GroundItems.Range(func(key, value interface{}) bool {
		giLength++
		return true
	})
	var pLength = 0
	r.Players.Range(func(key, value interface{}) bool {
		pLength++
		return true
	})
	if pLength == 0 && giLength == 0 {
		r.MarkedForDeletion = true
	}
}

func (r *Region) OnEnter(player *Player) {
	r.MarkedForDeletion = false // race condition - incase player leaves and enters in a single tick

	// We must add the player to all 9 regions entered on change
	// so that they get updates about regions around themselves
	r.Players.Store(player.Id, player)
	r.GroundItems.Range(func(key, value interface{}) bool {
		g := value.(*GroundItem)
		if g.Owner != nil && g.Owner.Name != player.Name {
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

	r.WorldObjects.Range(func(key, value interface{}) bool {
		obj := value.(model.WorldObjectInterface)
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: obj.GetObjectId(),
			Position: obj.GetPosition(),
			Face:     0,
			Typ:      10,
			Player:   player,
		})
		return true
	})
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
	r.Players.Delete(player.Id)
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
		Position:   p,
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
			r.Players.Range(func(key, value interface{}) bool {
				player := value.(*Player)
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
					Position: g.Position,
					Player:   player,
					ItemId:   g.ItemId,
				})
				return true
			})
			return false
		}
		return true
	})
}

func (r *Region) SetWorldObject(worldObject model.WorldObjectInterface) {
	key := fmt.Sprintf("%d-%d", worldObject.GetPosition().X, worldObject.GetPosition().Y)
	r.WorldObjects.Store(key, worldObject)
	r.Players.Range(func(key, value interface{}) bool {
		player := value.(*Player)
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendObjectPacket{
			ObjectId: worldObject.GetObjectId(),
			Position: worldObject.GetPosition(),
			Face:     0,
			Typ:      10,
			Player:   player,
		})
		return true
	})
}

func (r *Region) GetWorldObject(position *model.Position) model.WorldObjectInterface {
	key := fmt.Sprintf("%d-%d", position.X, position.Y)
	obj, ok := r.WorldObjects.Load(key)
	if !ok {
		return nil
	}
	return obj.(model.WorldObjectInterface)
}

func (r *Region) GetPlayersAsInterface() []model.PlayerInterface {
	var players = make([]model.PlayerInterface, 0)
	r.Players.Range(func(key, value interface{}) bool {
		player := value.(*Player)
		players = append(players, player)
		return true
	})
	return players
}

func (r *Region) GetNpcsAsInterface() []model.NpcInterface {
	var npcs = make([]model.NpcInterface, 0)
	r.Npcs.Range(func(key, value interface{}) bool {
		npc := value.(*Npc)
		npcs = append(npcs, npc)
		return true
	})
	return npcs
}

func GetRegionIdByPosition(p *model.Position) uint16 {
	regionX := p.X >> 3
	regionY := p.Y >> 3
	return (regionX / 8 << 8) + (regionY / 8)
}
