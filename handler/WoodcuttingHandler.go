package handler

import (
	"fmt"
	"math/rand"
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet/outgoing"
)

type WoodcuttingAxe struct {
	Id          int
	AnimationId int
	Strength    int // the axes bonus to chop speed
}

var axes = []*WoodcuttingAxe{
	{Id: 6739, AnimationId: 2846, Strength: 10}, // dragon
	{Id: 1359, AnimationId: 867, Strength: 8},   // rune
	{Id: 1357, AnimationId: 867, Strength: 6},   // addy
	{Id: 1355, AnimationId: 871, Strength: 4},   // mith
	{Id: 1353, AnimationId: 875, Strength: 3},   // steel
	{Id: 1349, AnimationId: 877, Strength: 2},   // iron
	{Id: 1351, AnimationId: 879, Strength: 1},   // bronze
}

type WoodcuttingTree struct {
	LevelRequirement int
	LifePoints       int
	LogId            int
	StumpId          int
	Experience       int
	Respawn          int
	Defence          int // the trees toughness to chop
}

var NormalTree = &WoodcuttingTree{
	LevelRequirement: 1,
	LifePoints:       1,
	LogId:            1511,
	StumpId:          1342,
	Experience:       100,
	Respawn:          20,
	Defence:          0,
}
var OakTree = &WoodcuttingTree{
	LevelRequirement: 15,
	LifePoints:       6,
	LogId:            1521,
	StumpId:          1341,
	Experience:       60,
	Respawn:          30,
	Defence:          2,
}
var WillowTree = &WoodcuttingTree{
	LevelRequirement: 30,
	LifePoints:       15,
	LogId:            1519,
	StumpId:          7399,
	Experience:       80,
	Respawn:          40,
	Defence:          3,
}
var MapleTree = &WoodcuttingTree{
	LevelRequirement: 45,
	LifePoints:       15,
	LogId:            1517,
	StumpId:          1343,
	Experience:       120,
	Respawn:          60,
	Defence:          5,
}
var YewTree = &WoodcuttingTree{
	LevelRequirement: 60,
	LifePoints:       20,
	LogId:            1515,
	StumpId:          7402,
	Experience:       160,
	Respawn:          80,
	Defence:          7,
}
var MagicTree = &WoodcuttingTree{
	LevelRequirement: 75,
	LifePoints:       20,
	LogId:            1513,
	StumpId:          7399,
	Experience:       205,
	Respawn:          10,
	Defence:          9,
}
var trees = map[int]*WoodcuttingTree{
	1315: NormalTree,
	1316: NormalTree,
	1318: NormalTree,
	1319: NormalTree,
	3033: NormalTree,
	1278: NormalTree,
	1276: NormalTree,
	1286: NormalTree,
	1282: NormalTree,
	1383: NormalTree,
	1281: OakTree,
	3037: OakTree,
	1308: WillowTree,
	5551: WillowTree,
	5552: WillowTree,
	5553: WillowTree,
	1307: MapleTree,
	4974: MapleTree,
	1309: YewTree,
	1306: MagicTree,
}

func ObjectIsWoodcuttingTree(objectId int) bool {
	if trees[objectId] != nil {
		return true
	}
	return false
}

type WoodcuttingHandler struct {
	player       *entity.Player
	treePosition *model.Position
	axe          *WoodcuttingAxe
	tree         *WoodcuttingTree
	treeId       int
	tickCount    int
}

// TODO: Check level requirement for tree and axe
func StartWoodcutting(treeId int, treePosition *model.Position, player *entity.Player) {
	weapon := player.Equipment.Items[outgoing.EQUIPMENT_SLOTS["weapon"]]
	for _, a := range axes {
		item := player.Inventory.FindItem(a.Id)
		if item != nil || weapon.ItemId == a.Id {
			treeWorldObject := entity.WorldProvider().GetWorldObject(treePosition)
			if treeWorldObject == nil {
				treeWorldObject = &TreeWorldObject{
					ObjectId:      treeId,
					Position:      treePosition,
					TreeId:        treeId,
					StumpId:       trees[treeId].StumpId,
					Respawn:       trees[treeId].Respawn,
					lifepoints:    trees[treeId].LifePoints,
					tickCount:     0,
					shouldRefresh: false,
				}
				entity.WorldProvider().SetWorldObject(treeWorldObject)
			}
			player.UpdateFlag.SetAnimation(a.AnimationId, 4)
			player.OngoingAction = &WoodcuttingHandler{
				axe:          a,
				tree:         trees[treeId],
				treePosition: treePosition,
				player:       player,
				treeId:       treeId,
			}
			player.UpdateFlag.SetFacePosition(treePosition)
			return
		}
	}
	player.OutgoingQueue = append(player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "You do not have an axe that you have the woodcutting level to use."})
}

func (w *WoodcuttingHandler) Tick() {
	if w.player.GetUpdateFlag().AnimationDuration <= 0 {
		w.player.UpdateFlag.SetAnimation(w.axe.AnimationId, 4)
	}

	if w.player.Inventory.IsFull() {
		w.player.OutgoingQueue = append(w.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: "Your inventory is too full to hold anymore logs."})
		w.StopWoodcutting()
		return
	}

	level := w.player.SkillHelper.Skills[model.Woodcutting].Level
	requirement := w.tree.LevelRequirement
	if level < requirement {
		w.player.OutgoingQueue = append(w.player.OutgoingQueue, &outgoing.SendMessagePacket{Message: fmt.Sprintf("You require level %d woodcutting to chop this tree.", w.tree.LevelRequirement)})
		w.StopWoodcutting()
		return
	}

	treeWorldObject := entity.WorldProvider().GetWorldObject(w.treePosition)
	if treeWorldObject == nil {
		w.StopWoodcutting()
		return
	}
	if treeWorldObject.GetObjectId() == w.tree.StumpId {
		w.StopWoodcutting()
		return
	}

	w.tickCount++
	if w.tickCount >= 2 {
		luck := 1 + rand.Float64()*(10-1)
		strength := w.axe.Strength
		defence := w.tree.Defence
		chance := float64(level-w.tree.LevelRequirement)/4 + float64(strength-defence) + 1
		if chance < 1.5 {
			chance = 1.5
		}
		if chance > 9.5 {
			chance = 9.5
		}
		//log.Printf("luck %+v, chance %+v", luck, chance)
		if chance >= luck {
			if t, ok := treeWorldObject.(*TreeWorldObject); ok {
				t.lifepoints--
				if t.lifepoints <= 0 {
					w.CutDownTree()
				}

				err := w.player.Inventory.AddItem(w.tree.LogId, 1)
				if err != nil {
					w.StopWoodcutting()
				}
				w.player.SkillHelper.AddExperience(model.Woodcutting, w.tree.Experience)
			}
		}

		w.tickCount = 0
	}
}

type TreeWorldObject struct {
	ObjectId      int
	Position      *model.Position
	TreeId        int
	StumpId       int
	Respawn       int
	lifepoints    int
	tickCount     int
	shouldRefresh bool
}

func (t *TreeWorldObject) GetPosition() *model.Position {
	return t.Position
}

func (t *TreeWorldObject) GetObjectId() int {
	return t.ObjectId
}

// TODO: Should only show the world object when it's cut down and delete it when it respawns..
func (t *TreeWorldObject) Tick() {
	t.tickCount++
	if t.tickCount >= t.Respawn {
		t.ObjectId = t.TreeId
		t.shouldRefresh = true
	}
}

func (t *TreeWorldObject) ShouldRefresh() bool {
	return t.shouldRefresh
}

func (w *WoodcuttingHandler) CutDownTree() {
	entity.WorldProvider().SetWorldObject(&TreeWorldObject{
		ObjectId:  w.tree.StumpId,
		Position:  w.treePosition,
		TreeId:    w.treeId,
		StumpId:   w.tree.StumpId,
		Respawn:   w.tree.Respawn,
		tickCount: 0,
	})
	w.StopWoodcutting()
}

func (w *WoodcuttingHandler) StopWoodcutting() {
	w.player.UpdateFlag.ClearAnimation()
	w.player.OngoingAction = nil
}
