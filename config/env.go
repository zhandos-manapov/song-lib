package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT   string
	GO_ENV string

	// Database variables
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
	DATABASE_HOST     string
	DATABASE_PORT     string

	// External api url
	API_URL string
}

func NewEnv(filename string) Config {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return Config{
		PORT:              getEnv("PORT"),
		GO_ENV:            getEnv("GO_ENV"),
		DATABASE_USER:     getEnv("DATABASE_USER"),
		DATABASE_PASSWORD: getEnv("DATABASE_PASSWORD"),
		DATABASE_NAME:     getEnv("DATABASE_NAME"),
		DATABASE_HOST:     getEnv("DATABASE_HOST"),
		DATABASE_PORT:     getEnv("DATABASE_PORT"),
		API_URL:           getEnv("API_URL"),
	}
}

func getEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		errorMessage := fmt.Sprintf("Env variable %s doesn't exists", key)
		log.Fatal(errorMessage)
		return ""
	}
}
