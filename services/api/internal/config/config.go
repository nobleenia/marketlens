package config

import (
	"fmt"
	"os"
	"strconv"
)

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}

type Config struct {
	Env  string
	Port string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	RedisHost string
	RedisPort string

	AfricaTalkingUser   string
	AfricaTalkingAPIKey string

	USSDWebHookSecret string

	AdminAPIKey   string
	USSDRateLimit int // requests per minute per IP
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

		AfricaTalkingUser:   getEnv("AT_USER", "marketlens"),
		AfricaTalkingAPIKey: getEnv("AT_API", "africatalking_api_key"),

		USSDWebHookSecret: getEnv("USSD_WEBHOOK_SECRET", ""),

		AdminAPIKey:   getEnv("ADMIN_API_KEY", ""),
		USSDRateLimit: getEnvInt("USSD_RATE_LIMIT", 30),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}
