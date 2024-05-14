package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log/slog"
	"nphud/pkg/config"
	"nphud/pkg/store"
	"os"
)

const migrationsDir = "file://cmd/migrate/migrations"

var ErrUnknownCommand = errors.New("unknown command")

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func run(args []string, getenv func(string) string) error {
	app := config.NewAppConfig(getenv)
	cmd := args[0]
	database, err := store.GetOrCreate(app.Database.FullPath)
	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(database, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsDir, app.Database.Name, driver)
	if err != nil {
		return err
	}

	slog.Info("starting migrations", "cmd", cmd)

	switch cmd {
	case "up":
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	case "down":
		if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	default:
		slog.Error("unknown command", "got", cmd, "expected", "up, down")
		return ErrUnknownCommand
	}

	slog.Info("migration complete")
	return nil
}

func main() {
	cmd := os.Args[len(os.Args)-1]
	args := []string{cmd}

	if err := run(args, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
