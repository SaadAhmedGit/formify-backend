package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SaadAhmedGit/forms/internal/config"
	"github.com/SaadAhmedGit/forms/internal/models"
)

var (
	db, err = createConnection()
)

func NewDatabase() (*gorm.DB, error) {
	return db, nil
}

func createConnection() (*gorm.DB, error) {
	env, err := config.Env()

	db, err := gorm.Open(sqlite.Open(env.DB_DSN), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database.")
		return db, err
	}

	db.AutoMigrate(&models.User{})

	return db, nil
}
