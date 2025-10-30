package config

import "os"

type Config struct {
	DSN  string
	Port string
}

func Load() Config {
	return Config{
		DSN:  getenv("DB_DSN", "host=localhost user=app password=app dbname=appdb port=5432 sslmode=disable TimeZone=UTC"),
		Port: getenv("PORT", "8080"),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
