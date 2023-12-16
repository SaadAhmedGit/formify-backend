package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

var (
	cfg, err = createEnv()
)

type config struct {
	// General
	APP_NAME string `env:"APP_NAME" envDefault:"Forms"`
	PORT     int    `env:"PORT" envDefault:"3000"`

	// Database
	DB_CONN_STRING string `env:"DB_CONN_STRING"`
	DB_HOST        string `env:"DB_HOST" envDefault:"localhost"`
	DB_PORT        int    `env:"DB_PORT" envDefault:"5432"`
	DB_USER        string `env:"DB_USER" envDefault:"postgres"`
	DB_PASSWORD    string `env:"DB_PASSWORD" envDefault:"postgres"`
	DB_NAME        string `env:"DB_NAME" envDefault:"forms"`

	// Dev
	IS_DEV          bool   `env:"IS_DEV" envDefault:"true"`
	DEV_CLIENT_URL  string `env:"DEV_CLIENT_URL" envDefault:"http://localhost:8080"`
	DEV_SERVER_URL  string `env:"DEV_SERVER_URL" envDefault:"http://localhost:3000"`
	PROD_CLIENT_URL string `env:"PROD_CLIENT_URL" envDefault:"https://forms-frontend.vercel.app"`
	PROD_SERVER_URL string `env:"PROD_SERVER_URL" envDefault:"https://forms-backend.vercel.app"`

	// Sendgrid email API
	SENDGRID_API_KEY string `env:"SENDGRID_API_KEY"`
	EMAIL_FROM       string `env:"EMAIL_FROM"`

	// JWT
	JWT_SECRET                     string `env:"JWT_SECRET"`
	JWT_ACCOUNT_ACTIVATION         string `env:"JWT_ACCOUNT_ACTIVATION"`
	JWT_RESET_PASSWORD             string `env:"JWT_RESET_PASSWORD"`
	TOKEN_VALIDITY_DAYS            int    `env:"TOKEN_VALIDITY_DAYS" envDefault:"7"`
	PRESIGNUP_TOKEN_VALIDITY_HOURS int    `env:"PRESIGNUP_TOKEN_VALIDITY_HOURS" envDefault:"5"`

	// Test environment variables
	TEST_DB_HOST     string `env:"TEST_DB_HOST" envDefault:"localhost"`
	TEST_DB_PORT     int    `env:"TEST_DB_PORT" envDefault:"5432"`
	TEST_DB_CONN_STR string `env:"TEST_DB_CONN_STR"`
	TEST_DB_USER     string `env:"TEST_DB_USER" envDefault:"postgres"`
	TEST_DB_PASSWORD string `env:"TEST_DB_PASSWORD" envDefault:"postgres"`
	TEST_DB_NAME     string `env:"TEST_DB_NAME" envDefault:"forms_test"`
}

func Env() (config, error) {
	return cfg, err
}

func createEnv() (config, error) {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("Failed to parse environment variables.")
		return cfg, err
	}
	return cfg, nil
}
