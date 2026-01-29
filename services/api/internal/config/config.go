package config

import "os"

type Config struct {
	Env  string
	Port string

	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string

	RedisHost string
	RedisPort string
}

func FromEnv() Config {
	return Config{
		Env:  getEnv("API_ENV", "local"),
		Port: getEnv("API_PORT", "8080"),

		PostgresHost:     getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnv("POSTGRES_DB", "marketlens"),
		PostgresUser:     getEnv("POSTGRES_USER", "marketlens"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "marketlens"),

		RedisHost: getEnv("REDIS_HOST", "redis"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
