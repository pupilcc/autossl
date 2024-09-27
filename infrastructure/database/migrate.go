package database

import (
	"context"
	"github.com/labstack/gommon/log"
)

func AutoMigrate() {
	client := initDB()
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
