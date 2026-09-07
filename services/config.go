package services

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	MailgunAPIKey string `env:"MAILGUN_API_KEY"`
	MailgunDomain string `env:"MAILGUN_DOMAIN"`
	JWTSecret     string `env:"JWT_SECRET"`
}

func LoadEnv() (Config, error) {
	cfg := Config{}
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
