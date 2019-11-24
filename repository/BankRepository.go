package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"rsps/model"
)

var NoBankItemsFoundError = errors.New("no inventory items found")

var BankSchema = `
	CREATE TABLE bank (
		playerName varchar(100) not null,
		slot int not null,
		itemId int not null,
		amount int not null,
	PRIMARY KEY (playerName, slot)
)`

type BankRepository struct {
	db *sqlx.DB
}

var BankRepositorySingleton *BankRepository
func NewBankRepository(db *sqlx.DB) {
	BankRepositorySingleton = &BankRepository{db: db}
}

func (b *BankRepository) Save(playerName string, items []*model.Item) {
	tx, err := b.db.Begin()
	if err != nil {
		return
	}

	for k, v := range items {
		_, _ = tx.Exec("INSERT INTO bank (playerName, slot, itemId, amount) values (?, ?, ?, ?) ON DUPLICATE KEY UPDATE itemId=?, amount=?",
			playerName,
			k,
			v.ItemId,
			v.Amount,
			v.ItemId,
			v.Amount)
	}

	tx.Commit()
}

func (b *BankRepository) Load(playerName string) ([]*model.Item, error) {
	rows, err := b.db.Query("SELECT slot, itemId, amount FROM bank WHERE playerName = ?", playerName)
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
		return nil, NoBankItemsFoundError
	}

	itemList := make([]*model.Item, len(itemMap))
	for k, v := range itemMap {
		itemList[k] = v
	}

	return itemList, nil
}
