package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"rsps/model"
)

var NoInventoryItemsFoundError = errors.New("no inventory items found")

var InventorySchema = `
	CREATE TABLE inventory (
		playerName varchar(100) not null,
		slot int not null,
		itemId int not null,
		amount int not null,
	PRIMARY KEY (playerName, slot)
)`

type InventoryRepository struct {
	db *sqlx.DB
}

var InventoryRepositorySingleton *InventoryRepository
func NewInventoryRepository(db *sqlx.DB) {
	InventoryRepositorySingleton = &InventoryRepository{db:db}
}

func (i *InventoryRepository) Save(playerName string, items []*model.Item) {
	tx, err := i.db.Begin()
	if err != nil {
		return
	}

	for k, v := range items {
		_, _ = tx.Exec("INSERT INTO inventory (playerName, slot, itemId, amount) values (?, ?, ?, ?) ON DUPLICATE KEY UPDATE itemId=?, amount=?",
			playerName,
			k,
			v.ItemId,
			v.Amount,
			v.ItemId,
			v.Amount)
	}

	tx.Commit()
}

func (i *InventoryRepository) Load(playerName string) ([]*model.Item, error) {
	rows, err := i.db.Query("SELECT slot, itemId, amount FROM inventory WHERE playerName = ?", playerName)
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
		return nil, NoInventoryItemsFoundError
	}

	itemList := make([]*model.Item, len(itemMap))
	for k, v := range itemMap {
		itemList[k] = v
	}

	return itemList, nil
}
