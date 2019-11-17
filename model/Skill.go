package model

import "time"

type Skill struct {
	Id         int
	Level      int
	Experience int
	BoostTimer time.Time
}

type SkillId int
const (
	Attack SkillId = iota
	Defence
	Strength
	Hitpoints
	Ranged
	Prayer
	Magic
	Cooking
	Woodcutting
	Fletching
	Fishing
	Firemaking
	Crafting
	Smithing
	Mining
	Hreblore
	Agility
	Thieving
	Slayer
	Farming
	Runecrafting
)
