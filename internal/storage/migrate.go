package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// Migrate runs the database migrations
func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return (err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		os.Getenv("DB_NAME"),
		driver)
	if err != nil {
		return (err)
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return (err)
	}

	fmt.Println("Database migrated successfully")
	return nil
}
