package main

import (
	"context"
	"log"
	"os"
	"song-lib/config"
	"song-lib/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	env := config.NewEnv(".env")
	context := context.Background()

	dbConfig := db.DbConfig{
		User:     env.DATABASE_USER,
		Password: env.DATABASE_PASSWORD,
		Host:     env.DATABASE_HOST,
		Port:     env.DATABASE_PORT,
		Name:     env.DATABASE_NAME,
	}

	db, err := db.NewDatabase(context, dbConfig).Connect()
  if err != nil {
    log.Fatal(err.Error())
  }
	defer db.Disconnect()

	sqlDb := stdlib.OpenDBFromPool(db.GetPool())
	defer sqlDb.Close()

	driver, err := pgx.WithInstance(sqlDb, &pgx.Config{})
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

	cmd := os.Args[len(os.Args)-1]
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