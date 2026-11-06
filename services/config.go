package services

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret     string `env:"JWT_SECRET"`
	STMP_SERVER   string `env:"STMP_SERVER"`
	STMP_DOMAIN   string `env:"STMP_DOMAIN"`
	STMP_PASSWORD string `env:"STMP_PASSWORD"`
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
