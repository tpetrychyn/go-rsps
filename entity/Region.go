package entity

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"rsps/model"
	"rsps/net/packet/outgoing"
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
	Players           map[uuid.UUID]*Player
	GroundItems       map[uuid.UUID]*GroundItem
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
	groundItems := make(map[uuid.UUID]*GroundItem, 0)
	if id == 12850 {
		groundItems[uuid.New()] = &GroundItem{
			Position: &model.Position{
				X: 3200,
				Y: 3200,
				Z: 0,
			},
			ItemId:    1351,
			Amount:    1,
			Owner:     nil,
			CreatedAt: time.Now(),
		}
	}

	return &Region{
		Id:          id,
		GroundItems: groundItems,
		Players:     make(map[uuid.UUID]*Player),
		WorldObjects: make(map[string]model.WorldObjectInterface),
	}
}

func (r *Region) Tick() {
	//for k, v := range r.GroundItems {
	//	if time.Now().Sub(v.CreatedAt) > 5*time.Second {
	//		r.RemovedItems = append(r.RemovedItems, v)
	//		delete(r.GroundItems, k)
	//	}
	//}

	for _, item := range r.GroundItems {
		if item.Owner != nil && time.Now().Sub(item.CreatedAt) > 5*time.Second {
			for _, player := range r.Players {
				if player.Id == item.Owner.Id {
					continue
				}
				player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
					Position:   item.Position,
					Player:     player,
					ItemId:     item.ItemId,
					ItemAmount: item.Amount,
				})
			}
			item.Owner = nil
		}
	}

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

	if len(r.Players) == 0 && len(r.GroundItems) == 0 {
		r.MarkedForDeletion = true
	}
}

func (r *Region) OnEnter(player *Player) {
	for _, v := range r.GroundItems {
		if v.Owner != nil && v.Owner != player {
			continue
		} // only show it if nobody owns it or you own it
		player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.CreateGroundItemPacket{
			Position:   v.Position,
			Player:     player,
			ItemId:     v.ItemId,
			ItemAmount: v.Amount,
		})
	}

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
	top := r.Id - 1
	bottom := r.Id + 1
	// includes r itself
	return []uint16{top - 256, top, top + 256, r.Id - 256, r.Id, r.Id + 256, bottom - 256, bottom, bottom + 256}
}

func (r *Region) FindGroundItemByPosition(id int, p *model.Position) *GroundItem {
	for _, v := range r.GroundItems {
		if v.ItemId == id && v.Position.X == p.X && v.Position.Y == p.Y {
			return v
		}
	}
	return nil
}

func (r *Region) CreateGroundItemAtPosition(dropper *Player, item *model.Item, p *model.Position) {
	r.GroundItems[uuid.New()] = &GroundItem{
		Position:  p,
		ItemId:    item.ItemId,
		Amount:    item.Amount,
		Owner:     dropper,
		CreatedAt: time.Now(),
	}

	dropper.OutgoingQueue = append(dropper.OutgoingQueue, &outgoing.CreateGroundItemPacket{
		Position:   dropper.Position,
		Player:     dropper,
		ItemId:     item.ItemId,
		ItemAmount: item.Amount,
	})
}

func (r *Region) RemoveGroundItemIdAtPosition(id int, position *model.Position) {
	for k, v := range r.GroundItems {
		if v.ItemId == id && v.Position.X == position.X && v.Position.Y == position.Y {
			delete(r.GroundItems, k)
			for _, p := range r.Players {
				p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.RemoveGroundItemPacket{
					Position: v.Position,
					Player:   p,
					ItemId:   v.ItemId,
				})
			}
			return
		}
	}
	log.Printf("item not found to remove")
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
