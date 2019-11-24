package entity

import (
	"fmt"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/repository"
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
			Message:     fmt.Sprintf("Congratulations, you just advanced a %s level!", LevelupInterfaces[skillId].SkillName),
		})
		s.Player.OutgoingQueue = append(s.Player.OutgoingQueue, &outgoing.SendTextInterfacePacket{
			InterfaceId: LevelupInterfaces[skillId].SecondLine,
			Message:     fmt.Sprintf("Your %s level is now %d.", LevelupInterfaces[skillId].SkillName, s.Skills[skillId].Level),
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

	go repository.SkillRepositorySingleton.Save(s.Player.Name, s.Skills)
}

type LevelupInterface struct {
	InterfaceId uint
	FirstLine   uint
	SecondLine  uint
	SkillName   string
}

var LevelupInterfaces = map[model.SkillId]LevelupInterface{
	model.Attack: {
		InterfaceId: 6247,
		FirstLine:   6248,
		SecondLine:  6249,
		SkillName:   "attack",
	},
	model.Woodcutting: {
		InterfaceId: 4272,
		FirstLine:   4273,
		SecondLine:  4274,
		SkillName:   "woodcutting",
	},
	model.Thieving: {
		InterfaceId: 4261,
		FirstLine:   4263,
		SecondLine:  4264,
		SkillName:   "thieving",
	},
}
