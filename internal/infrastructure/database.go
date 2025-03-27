package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresDB() *sql.DB {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", err)
	}

	err = runMigrations(dsn)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	return db
}

func runMigrations(dsn string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsPath := filepath.Join(wd, "internal/infrastructure/migrations")
	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
