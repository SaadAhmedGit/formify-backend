package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

var (
	cfg, err = createEnv()
)

type config struct {
	APP_NAME               string `env:"APP_NAME" envDefault:"Forms"`
	PORT                   int    `env:"PORT" envDefault:"3000"`
	DB_DSN                 string `env:"DB_DSN" envDefault:"forms.sqlite"`
	IS_DEV                 bool   `env:"IS_DEV" envDefault:"true"`
	DEV_CLIENT_URL         string `env:"DEV_CLIENT_URL" envDefault:"http://localhost:8080"`
	DEV_SERVER_URL         string `env:"DEV_SERVER_URL" envDefault:"http://localhost:3000"`
	PROD_CLIENT_URL        string `env:"PROD_CLIENT_URL" envDefault:"https://forms-frontend.vercel.app"`
	PROD_SERVER_URL        string `env:"PROD_SERVER_URL" envDefault:"https://forms-backend.vercel.app"`
	SENDGRID_API_KEY       string `env:"SENDGRID_API_KEY"`
	EMAIL_FROM             string `env:"EMAIL_FROM"`
	JWT_SECRET             string `env:"JWT_SECRET"`
	JWT_ACCOUNT_ACTIVATION string `env:"JWT_ACCOUNT_ACTIVATION"`
	JWT_RESET_PASSWORD     string `env:"JWT_RESET_PASSWORD"`
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
