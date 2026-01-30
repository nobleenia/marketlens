package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"marketlens/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		dir    = flag.String("dir", "migrations", "migrations directory (relative to services/api)")
		action = flag.String("action", "status", "goose action: up|down|status|redo|version")
	)
	flag.Parse()
	cfg := config.FromEnv()

	// Use pgx stdlib driver for database/sql compatibility
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Ensure we use Postgres dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	// goose expects migrations dir relative to current working directory
	// We'll run this binary from services/api by Makefile.
	if err := runGoose(db, *dir, *action); err != nil {
		log.Fatalf("goose %s migration failed: %v", *action, err)
	}
}

func runGoose(db *sql.DB, dir, action string) error {
	switch action {
	case "up":
		return goose.Up(db, dir)
	case "down":
		return goose.Down(db, dir)
	case "status":
		return goose.Status(db, dir)
	case "redo":
		return goose.Redo(db, dir)
	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			return err
		}
		_, _ = fmt.Printf("current db version: %d\n", version)
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}
