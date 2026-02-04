package config

import (
	"os"
	"testing"
)

func TestPostgresDSN_FormatsCorrectly(t *testing.T) {
	cfg := Config{
		PostgresUser:     "testuser",
		PostgresPassword: "testpass",
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresDB:       "testdb",
	}

	expectedDSN := "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if dsn := cfg.PostgresDSN(); dsn != expectedDSN {
		t.Errorf("expected DSN %q, got %q", expectedDSN, dsn)
	}
}

func TestFromEnv_UsesDefaults(t *testing.T) {
	// Clear environment variables for the test
	os.Unsetenv("API_ENV")
	os.Unsetenv("API_PORT")
	os.Unsetenv("POSTGRES_HOST")
	os.Unsetenv("POSTGRES_PORT")
	os.Unsetenv("POSTGRES_DB")
	os.Clearenv()

	cfg := FromEnv()

	if cfg.Env != "local" {
		t.Errorf("expected Env to be 'local', got %q", cfg.Env)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port to be '8080', got %q", cfg.Port)
	}
	if cfg.PostgresHost != "postgres" {
		t.Errorf("expected PostgresHost to be 'postgres', got %q", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "5432" {
		t.Errorf("expected PostgresPort to be '5432', got %q", cfg.PostgresPort)
	}
	if cfg.PostgresDB != "marketlens" {
		t.Errorf("expected PostgresDB to be 'marketlens', got %q", cfg.PostgresDB)
	}
}

func TestFromEnv_UsesEnvironmentVariables(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("API_ENV", "production")
	os.Setenv("API_PORT", "9090")
	os.Setenv("POSTGRES_HOST", "db.example.com")
	os.Setenv("POSTGRES_PORT", "6543")
	os.Setenv("POSTGRES_DB", "prod_db")
	defer func() {
		os.Unsetenv("API_ENV")
		os.Unsetenv("API_PORT")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_DB")
	}()

	cfg := FromEnv()

	if cfg.Env != "production" {
		t.Errorf("expected Env to be 'production', got %q", cfg.Env)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected Port to be '9090', got %q", cfg.Port)
	}
	if cfg.PostgresHost != "db.example.com" {
		t.Errorf("expected PostgresHost to be 'db.example.com', got %q", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "6543" {
		t.Errorf("expected PostgresPort to be '6543', got %q", cfg.PostgresPort)
	}
	if cfg.PostgresDB != "prod_db" {
		t.Errorf("expected PostgresDB to be 'prod_db', got %q", cfg.PostgresDB)
	}
}

func TestGetEnv_ReturnsDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_VAR")

	got := getEnv("NONEXISTENT_VAR", "default_value")
	if got != "default_value" {
		t.Errorf("expected default value 'default_value', got %q", got)
	}
}

func TestGetEnv_ReturnsEnvValue(t *testing.T) {
	os.Setenv("EXISTENT_VAR", "env_value")
	defer os.Unsetenv("EXISTENT_VAR")

	got := getEnv("EXISTENT_VAR", "default_value")
	if got != "env_value" {
		t.Errorf("expected env value 'env_value', got %q", got)
	}
}
