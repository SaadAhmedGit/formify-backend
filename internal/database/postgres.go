package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/SaadAhmedGit/forms/internal/config"
	"github.com/SaadAhmedGit/forms/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db, err = createConnection()
)

func NewDatabase() (*sqlx.DB, error) {
	return db, nil
}

func buildSchema() string {
	schemaArray := []string{}
	schemaArray = append(schemaArray, models.CREATE_USERS_TABLE_QUERY)

	schema := strings.Join(schemaArray, "\n\n")

	return schema
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

	schema := buildSchema()
	db.MustExec(schema)

	return db, nil
}
