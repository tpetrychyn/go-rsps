package model

import (
	"math"
	"time"
)

type Skill struct {
	Id         int
	Level      int
	Experience int
	BoostTimer time.Time
}

func (s *Skill) GetCurrentLevel() int {
	return s.Level
}

func (s *Skill) GetLevelForExperience() int {
	return getLevelForExperience(s.Experience)
}

func getLevelForExperience(experience int) int {
	points := 0.0
	if experience > 13034430 {
		return 99
	}
	for l := 1; l < 99; l++ {
		points += math.Floor(float64(l) + 300.0*math.Pow(2.0, float64(l)/7.0))
		if int(math.Floor(points/4)) >= experience {
			return l
		}
	}
	return 1
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
