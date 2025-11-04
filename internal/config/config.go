package config

import (
	env "github.com/caarlos0/env/v10"
)

type Config struct {
	Port         string `env:"PORT,required"`
	PostgresUser string `env:"POSTGRES_USER,required"`
	PostgresPass string `env:"POSTGRES_PASS,required"`
	PostgresHost string `env:"POSTGRES_HOST,required"`
	PostgresDb   string `env:"POSTGRES_DB,required"`
	RedisHost    string `env:"REDIS_HOST,required"`
	RedisPort    string `env:"REDIS_PORT,required"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
