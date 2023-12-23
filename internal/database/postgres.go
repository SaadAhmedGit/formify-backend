package database

import (
	"fmt"
	"log"

	"github.com/SaadAhmedGit/formify/internal/config"
	"github.com/SaadAhmedGit/formify/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db, _ = createConnection()
)

func NewDatabase() (*sqlx.DB, error) {
	return db, nil
}

func createTables(db *sqlx.DB) {
	table_creation_queries := []string{models.CREATE_FORMS_TABLE_QUERY, models.CREATE_USERS_TABLE_QUERY}

	for _, query := range table_creation_queries {
		db.Exec(query)
	}
}

func createConnection() (*sqlx.DB, error) {
	env, err := config.Env()

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		env.DB_HOST, env.DB_PORT, env.DB_USER, env.DB_PASSWORD, env.DB_NAME)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database. %s\n", err.Error())
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database.")
		return db, err
	}

	createTables(db)

	return db, nil
}
