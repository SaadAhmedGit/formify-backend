package tests

import (
	"fmt"
	"log"
	"strings"
	"time"

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

func buildSchemaQuery() string {
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

	// https://github.com/lib/pq/issues/835#issuecomment-1774313464
	var db *sqlx.DB
	for i := 0; i < 20; i++ {
		db, err = sqlx.Connect("postgres", connStr)
		if err == nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Printf("Failed to connect to database. %s\n", err.Error())
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database.")
		return db, err
	}

	schemaQuery := buildSchemaQuery()
	db.MustExec(schemaQuery)

	return db, nil
}
