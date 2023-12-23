package tests

import (
	"fmt"
	"log"
	"strings"

	"github.com/SaadAhmedGit/formify/internal/config"
	"github.com/SaadAhmedGit/formify/internal/models"

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
		env.TEST_DB_HOST, env.TEST_DB_PORT, env.TEST_DB_USER, env.TEST_DB_PASSWORD, env.TEST_DB_NAME)
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

	//schema := buildSchema()
	//db.MustExec(schema)

	return db, nil
}
