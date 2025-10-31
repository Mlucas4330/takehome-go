package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type Config struct {
	Port       int    `env:"PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPass     string `env:"DB_PASS,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBPort     int    `env:"DB_PORT,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
	DBTimezone string `env:"DB_TIMEZONE,required"`
}

func LoadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("it was not possible to load .env variables: %+v", err)
	}
	return &cfg
}
