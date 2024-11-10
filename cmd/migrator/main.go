package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
)

func main() {
	var storagePath, migrationsPath, migrationTable string

	flag.StringVar(&storagePath, "storagePath", "./storage", "Path to store files")
	flag.StringVar(&migrationsPath, "migrationsPath", "./migrations", "Path to store files")
	flag.StringVar(&migrationTable, "migrationTable", "./migrations", "Path to store files")
	flag.Parse()

	if storagePath == "" {
		panic("storage pass is required")
	}

	if migrationsPath == "" {
		panic("migrations pass is required")
	}

	m, err := migrate.New("file://"+migrationsPath,
		fmt.Sprintf("postgres://%s/%s", storagePath, migrationTable),
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No change")
			return
		}
		panic(err)
	}

	fmt.Println("Migration complete")
}
