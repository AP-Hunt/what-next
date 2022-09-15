package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

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

func InTransaction[T interface{}](
	callback func(tx *sqlx.Tx) (*T, error),
	db *sqlx.DB,
	ctx context.Context,
) (*T, error) {
	tx, err := db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("starting transaction: %s", err)
	}

	val, err := callback(tx)

	if err != nil {
		if e := tx.Rollback(); e != nil {
			return nil, fmt.Errorf("error rolling back transaction: %s", e)
		}

		return nil, fmt.Errorf("error in transaction: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error comitting transaction: %s", err)
	}

	return val, nil
}
