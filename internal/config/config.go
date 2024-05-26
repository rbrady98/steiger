package config

import "os"

type Config struct {
	Env   string
	Port  string
	DbURL string
}

func NewConfig() Config {
	return Config{
		Env:   getEnv("APP_ENV", "local"),
		Port:  getEnv("PORT", "8000"),
		DbURL: getEnv("DB_URL", "./test.db"),
	}
}

func getEnv(key string, fallback string) string {
	if v, set := os.LookupEnv(key); set {
		return v
	}

	return fallback
}
