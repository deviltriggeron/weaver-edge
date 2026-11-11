package app

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"weaver-edge/config"
	"weaver-edge/internal/infrastructure/db"
)

func Migrate(cfg config.Config) error {
	db, err := db.InitDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch command {
	case "up":
		if err := goose.Up(db, "migrations"); err != nil {
			return err
		}
	case "down":
		if err := goose.Down(db, "migrations"); err != nil {
			return err
		}
	case "status":
		if err := goose.Status(db, "migrations"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}
