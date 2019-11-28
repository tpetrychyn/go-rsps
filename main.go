package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"math/rand"
	"os"
	"rsps/handler"
	"rsps/net"
	"rsps/repository"
	"rsps/util"
	"time"
)

func main() {
	db := initDatabase()
	playerRepository := repository.NewPlayerRepository(db)

	util.LoadItemDefinitions()

	rand.Seed(time.Now().Unix())

	handler.LoadScripts()

	server := net.NewTcpServer(43594)
	server.Start(playerRepository)

	//util.LoadCache()
}

func initDatabase() *sqlx.DB {
	db, err := sqlx.Connect("mysql", os.Getenv("MARIADB_DSN"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(repository.PlayerSchema)
	if err != nil {
		log.Printf("playerechema %s", err.Error())
	}
	_, err = db.Exec(repository.InventorySchema)
	if err != nil {
		log.Printf("inventoryschema %s", err.Error())
	}
	_, err = db.Exec(repository.EquipmentSchema)
	if err != nil {
		log.Printf("equipmentschema %s", err.Error())
	}
	_, err = db.Exec(repository.BankSchema)
	if err != nil {
		log.Printf("bankschema %s", err.Error())
	}
	_, err = db.Exec(repository.SkillSchema)
	if err != nil {
		log.Printf("skillschema %s", err.Error())
	}

	repository.NewInventoryRepository(db)
	repository.NewEquipmentRepository(db)
	repository.NewBankRepository(db)
	repository.NewSkillRepository(db)
	return db
}
