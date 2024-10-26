package main

import (
	"context"
	"fmt"
	"log"
	"song-lib/config"
	"song-lib/db"
)

func startServer() {

	context := context.Background()
	env := config.NewEnv(".env")

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

	module := NewModule(context, db, env)

	router := NewRouter()
	router.LoadControllers("/api/v1", module.Controllers())

	log.Printf(`Server is running on port %s`, env.PORT)

	if err := router.Run(fmt.Sprintf(`:%s`, env.PORT)); err != nil {
		log.Fatal(err.Error())
	}
}
