package entity

import (
	"fmt"
	"math"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"time"
)

type SkillHelper struct {
	Player *Player
	Skills Skills
}

type Skills = []*model.Skill

func NewSkillHelper(p *Player) *SkillHelper {
	s := make([]*model.Skill, 20)
	for i:=0;i<20;i++ {
		s[i] = &model.Skill{
			Id:         i,
			Level:      1,
			Experience: 0,
			BoostTimer: time.Time{},
		}
	}

	return &SkillHelper{
		Player: p,
		Skills: s,
	}
}

func getLevelForExperience(experience int) int {
	points := 0.0
	if experience > 13034430 {
		return 99
	}
	for l:=1;l<99;l++ {
		points += math.Floor(float64(l) + 300.0 * math.Pow(2.0, float64(l) / 7.0))
		if int(math.Floor(points / 4)) >= experience {
			return l
		}
	}
	return 1
}

func (s *SkillHelper) SetLevelToExperience(skillId model.SkillId, experience int) {
	s.Skills[skillId].Level = getLevelForExperience(experience)
	s.Skills[skillId].Experience = experience
	s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SetSkillLevelPacket{
		SkillNum:   int(skillId),
		Level:      s.Skills[skillId].Level,
		Experience: s.Skills[skillId].Experience,
	})
}

func (s *SkillHelper) SetLevel(skillId model.SkillId, level int) {
	s.Skills[skillId].Level = level
	s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SetSkillLevelPacket{
		SkillNum:   int(skillId),
		Level:      s.Skills[skillId].Level,
		Experience: s.Skills[skillId].Experience,
	})
}

func (s *SkillHelper) SetExperience(skillId model.SkillId, experience int) {
	s.Skills[skillId].Experience = experience
	s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SetSkillLevelPacket{
		SkillNum:   int(skillId),
		Level:      s.Skills[skillId].Level,
		Experience: s.Skills[skillId].Experience,
	})
}

func (s *SkillHelper) AddExperience(skillId model.SkillId, experience int) {
	s.Skills[skillId].Experience += experience
	if getLevelForExperience(s.Skills[skillId].Experience) > s.Skills[skillId].Level {
		s.Skills[skillId].Level += 1
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendTextInterfacePacket{
			InterfaceId: 6248,
			Message:     "Congratulations, you just advanced an attack level!",
		})
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendTextInterfacePacket{
			InterfaceId: 6249,
			Message:     fmt.Sprintf("Your attack level is now %d.", s.Skills[skillId].Level),
		})
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendChatInterfacePacket{
			InterfaceId: LevelupInterfaces[skillId],
		})
	}
	s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SetSkillLevelPacket{
		SkillNum:   int(skillId),
		Level:      s.Skills[skillId].Level,
		Experience: s.Skills[skillId].Experience,
	})
}

var LevelupInterfaces = map[model.SkillId]uint{
	model.Attack: 6247,
	model.Woodcutting: 4272,
}