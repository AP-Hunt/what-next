package context

import (
	"context"
	"os"
	"path"

	"github.com/AP-Hunt/what-next/m/db"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	CFG_KEY_DATA_DIR = "WHAT_NEXT_DATA_DIR"
)

func CreateDefaultCommandContext(parentContext context.Context) (CommandContext, error) {
	ctx := NewCommandContext(parentContext)

	initViper()

	database, err := initDb(viper.GetString(CFG_KEY_DATA_DIR))
	if err != nil {
		return CommandContext{}, err
	}

	ctx = ctx.WithDatabaseConnection(database)

	return ctx, nil
}

func initViper() {
	var homeDir string
	if d, err := os.UserHomeDir(); err != nil {
		homeDir = d
		panic(err)
	}

	viper.SetDefault(CFG_KEY_DATA_DIR, homeDir)
	viper.BindEnv(CFG_KEY_DATA_DIR)
}

func initDb(dataDir string) (*sqlx.DB, error) {
	dbPath := path.Join(dataDir, "what-next.sqlite")
	conn, err := db.Connect(dbPath)

	if err != nil {
		return nil, err
	}

	err = db.Migrate(conn.DB)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
