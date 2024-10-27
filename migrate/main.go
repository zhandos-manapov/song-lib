package main

import (
	"context"
	"log"
	"os"
	"song-lib/config"
	"song-lib/db"
	"song-lib/migrate/migrate"

	"github.com/jackc/pgx/v5/stdlib"
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

	cmd := os.Args[len(os.Args)-1]
	migrate.RunMigration(sqlDb, cmd)
}
