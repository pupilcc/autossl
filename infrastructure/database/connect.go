package database

import (
	"autossl/infrastructure/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"time"
)

var (
	dbClient *ent.Client
	once     sync.Once
)

func GetDBClient() *ent.Client {
	once.Do(func() {
		dbClient = initDB()
	})
	return dbClient
}

func initDB() *ent.Client {
	drv, err := sql.Open(dialect.SQLite, "/root/data/autossl.db?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return ent.NewClient(ent.Driver(drv))
}
