package main

import (
	"context"
	"fmt"
	"github.com/swaggo/http-swagger"
	"log"
	"song-lib/config"
	"song-lib/db"
	_ "song-lib/docs"
)

//	@title			Song Library API
//	@version		1.0
//	@description	This is song library API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3001
//	@BasePath	/api/v1
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

	mux := router.GetMux()
	mux.Get("/swagger/*", httpSwagger.Handler(
		// Visit http://localhost:3001/swagger/index.html for swagger docs
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", env.PORT)),
	))

	log.Printf(`Server is running on port %s`, env.PORT)

	if err := router.Run(fmt.Sprintf(":%s", env.PORT)); err != nil {
		log.Fatal(err.Error())
	}
}
