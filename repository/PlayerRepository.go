package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"rsps/model"
)

var PlayerNotFoundError = errors.New("player not found")

var PlayerSchema = `
CREATE TABLE player (
	playerName varchar(100) not null,
	password varchar(1000) not null,
	x int not null,
	y int not null,
	PRIMARY KEY (playerName)
)
`

type PlayerRepository struct {
	db *sqlx.DB
}

func NewPlayerRepository(db *sqlx.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (p *PlayerRepository) Load(name string) (*model.Position, []byte, error) {
	rows, err := p.db.Query("SELECT password, x, y from player WHERE playerName = ?", name)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var hashedPassword string
		var playerPosition model.Position
		err := rows.Scan(&hashedPassword, &playerPosition.X, &playerPosition.Y)
		if err != nil {
			return nil, nil, err
		}
		return &playerPosition, []byte(hashedPassword), nil
	}

	return nil, nil, PlayerNotFoundError
}

func (p *PlayerRepository) Create(name string, password []byte, position *model.Position) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	res, err := p.db.Exec("INSERT INTO player (playerName, password, x, y) values (?, ?, ?, ?)",
		name,
		hash,
		position.X,
		position.Y)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return nil, err
	}
	log.Printf("%+v" ,res)
	return hash, err
}

func (p *PlayerRepository) Save(name string, position *model.Position) error {
	res, err := p.db.Exec("UPDATE player SET x = ?, y = ? WHERE playerName = ?",
		position.X,
		position.Y,
		name)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
