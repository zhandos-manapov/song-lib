package migrate

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(instance *sql.DB, cmd string) {
	driver, err := pgx.WithInstance(instance, &pgx.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrate/migrations",
		"pgx5",
		driver,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}
	}
}
