package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"rsps/model"
)

var NoSkillFoundError = errors.New("no skills found")

var SkillSchema = `
	CREATE TABLE skill (
		playerName varchar(100) not null,
		skill int not null,
		level int not null,
		experience int not null,
	PRIMARY KEY (playerName, skill)
)`

type SkillRepository struct {
	db *sqlx.DB
}

var SkillRepositorySingleton *SkillRepository

func NewSkillRepository(db *sqlx.DB) {
	SkillRepositorySingleton = &SkillRepository{db: db}
}

func (e *SkillRepository) Save(playerName string, skills []*model.Skill) {
	tx, err := e.db.Begin()
	if err != nil {
		return
	}

	for k, v := range skills {
		_, _ = tx.Exec("INSERT INTO skill (playerName, skill, level, experience) values (?, ?, ?, ?) ON DUPLICATE KEY UPDATE level=?, experience=?",
			playerName,
			k,
			v.Level,
			v.Experience,
			v.Level,
			v.Experience)
	}

	tx.Commit()
}

func (e *SkillRepository) Load(playerName string) ([]*model.Skill, error) {
	rows, err := e.db.Query("SELECT skill, level, experience FROM skill WHERE playerName = ?", playerName)
	if err != nil {
		return nil, err
	}

	skillMap := make(map[int]*model.Skill)
	for rows.Next() {
		var skillId int
		var skill model.Skill
		err := rows.Scan(&skillId, &skill.Level, &skill.Experience)
		if err != nil {
			continue
		}
		skillMap[skillId] = &skill
	}

	if len(skillMap) == 0 {
		return nil, NoSkillFoundError
	}

	skillList := make([]*model.Skill, len(skillMap))
	for k, v := range skillMap {
		skillList[k] = v
	}

	return skillList, nil
}
