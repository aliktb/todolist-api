package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ali/todolist-api/internal/config"
)

func main() {
	cfg := config.Load()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("resolve migrations path: %v", err)
	}

	m, err := migrate.New("file://"+migrationsPath, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch command {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "reset":
		err = m.Down()
	case "version":
		version, dirty, versionErr := m.Version()
		if errors.Is(versionErr, migrate.ErrNilVersion) {
			fmt.Println("version: none")
			return
		}
		if versionErr != nil {
			log.Fatalf("read version: %v", versionErr)
		}
		fmt.Printf("version: %d dirty: %t\n", version, dirty)
		return
	default:
		log.Fatalf("unknown command %q, expected one of: up, down, reset, version", command)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("no schema change")
		return
	}
	if err != nil {
		log.Fatalf("migration command %q failed: %v", command, err)
	}

	fmt.Printf("migration command %q completed\n", command)
}
