package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	User   string
	Dbname string
}

type Config struct {
	Database DatabaseConfig
}

func New() *Config {
	return &Config{
		Database: DatabaseConfig{
			User:   getEnv("DB_USER"),
			Dbname: getEnv("DB_NAME"),
		},
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return value
}
