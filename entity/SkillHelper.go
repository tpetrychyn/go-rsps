package entity

import (
	"fmt"
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
	for i := 0; i < 20; i++ {
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

func (s *SkillHelper) SetLevelToExperience(skillId model.SkillId, experience int) {
	s.Skills[skillId].Level = s.Skills[skillId].GetLevelForExperience()
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
	if s.Skills[skillId].GetLevelForExperience() > s.Skills[skillId].Level {
		s.Skills[skillId].Level += 1
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendTextInterfacePacket{
			InterfaceId: LevelupInterfaces[skillId].FirstLine,
			Message:     "Congratulations, you just advanced an attack level!",
		})
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendTextInterfacePacket{
			InterfaceId: LevelupInterfaces[skillId].SecondLine,
			Message:     fmt.Sprintf("Your attack level is now %d.", s.Skills[skillId].Level),
		})
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendChatInterfacePacket{
			InterfaceId: LevelupInterfaces[skillId].InterfaceId,
		})
	}
	s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SetSkillLevelPacket{
		SkillNum:   int(skillId),
		Level:      s.Skills[skillId].Level,
		Experience: s.Skills[skillId].Experience,
	})
}

type LevelupInterface struct {
	InterfaceId uint
	FirstLine   uint
	SecondLine  uint
}

var LevelupInterfaces = map[model.SkillId]LevelupInterface{
	model.Attack: {
		InterfaceId: 6247,
		FirstLine:   6248,
		SecondLine:  6249,
	},
	model.Woodcutting: {
		InterfaceId: 4272,
		FirstLine:   4273,
		SecondLine:  4274,
	},
}
