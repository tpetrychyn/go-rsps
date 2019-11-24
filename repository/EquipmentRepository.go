package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"rsps/model"
)

var NoEquipmentItemsFoundError = errors.New("no inventory items found")

var EquipmentSchema = `
	CREATE TABLE equipment (
		playerName varchar(100) not null,
		slot int not null,
		itemId int not null,
		amount int not null,
	PRIMARY KEY (playerName, slot)
)`

type EquipmentRepository struct {
	db *sqlx.DB
}

var EquipmentRepositorySingleton *EquipmentRepository
func NewEquipmentRepository(db *sqlx.DB) {
	EquipmentRepositorySingleton = &EquipmentRepository{db: db}
}

func (e *EquipmentRepository) Save(playerName string, items []*model.Item) {
	tx, err := e.db.Begin()
	if err != nil {
		return
	}

	for k, v := range items {
		_, _ = tx.Exec("INSERT INTO equipment (playerName, slot, itemId, amount) values (?, ?, ?, ?) ON DUPLICATE KEY UPDATE itemId=?, amount=?",
			playerName,
			k,
			v.ItemId,
			v.Amount,
			v.ItemId,
			v.Amount)
	}

	tx.Commit()
}

func (e *EquipmentRepository) Load(playerName string) ([]*model.Item, error) {
	rows, err := e.db.Query("SELECT slot, itemId, amount FROM equipment WHERE playerName = ?", playerName)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[int]*model.Item)
	for rows.Next() {
		var slot int
		var item model.Item
		err := rows.Scan(&slot, &item.ItemId, &item.Amount)
		if err != nil { continue }
		itemMap[slot] = &item
	}

	if len(itemMap) == 0 {
		return nil, NoEquipmentItemsFoundError
	}

	itemList := make([]*model.Item, len(itemMap))
	for k, v := range itemMap {
		itemList[k] = v
	}

	return itemList, nil
}
