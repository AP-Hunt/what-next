package db

import (
	"database/sql"
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pressly/goose/v3"
)

func Connect(connectionString string) (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", connectionString)
}

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func Migrate(database *sql.DB) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(embeddedMigrations)

	err := goose.SetDialect("sqlite")
	if err != nil {
		return err
	}

	err = goose.Up(database, "migrations")
	if err != nil {
		return err
	}

	return nil
}
