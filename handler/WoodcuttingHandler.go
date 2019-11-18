package handler

import (
	"log"
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
	Experience:       40,
	Respawn:          20,
	Defence:          1,
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
	Respawn:          100,
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
	tickCount    int
}

// TODO: Check level requirement for tree and axe
func StartWoodcutting(treeId int, treePosition *model.Position, player *entity.Player) {
	player.SkillHelper.SetLevel(model.Woodcutting, 99)
	weapon := player.Equipment.Items[outgoing.EQUIPMENT_SLOTS["weapon"]]
	for _, a := range axes {
		item := player.Inventory.FindItem(a.Id)
		if item != nil || weapon.ItemId == a.Id {
			player.UpdateFlag.SetAnimation(a.AnimationId, 4)
			player.OngoingAction = &WoodcuttingHandler{
				axe:          a,
				tree:         trees[treeId],
				treePosition: treePosition,
				player:       player,
			}
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
	}

	level := w.player.SkillHelper.Skills[model.Woodcutting].Level
	requirement := w.tree.LevelRequirement
	if level < requirement {
		w.StopWoodcutting()
	}

	w.tickCount++
	if w.tickCount >= 4 {
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
		log.Printf("luck %+v, chance %+v", luck, chance)
		if chance >= luck {
			err := w.player.Inventory.AddItem(w.tree.LogId, 1)
			if err != nil {
				w.StopWoodcutting()
			}
			w.player.SkillHelper.AddExperience(model.Woodcutting, w.tree.Experience*EXP_MULTIPLIER)
			// w.CutDownTree()
			// TODO: Make regions have world objects and instantiate the tree when clicked and
			//  broadcast update when cut down
		}

		w.tickCount = 0
	}
}

var EXP_MULTIPLIER = 100000

func (w *WoodcuttingHandler) CutDownTree() {
	w.player.OutgoingQueue = append(w.player.OutgoingQueue, &outgoing.SendObjectPacket{
		ObjectId: w.tree.StumpId,
		Position: w.treePosition,
		Face:     0,
		Typ:      10,
		Player:   w.player,
	})
	w.StopWoodcutting()
}

func (w *WoodcuttingHandler) StopWoodcutting() {
	w.player.UpdateFlag.ClearAnimation()
	w.player.OngoingAction = nil
}
