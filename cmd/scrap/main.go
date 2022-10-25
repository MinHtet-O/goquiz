package main

import (
	"github.com/joho/godotenv"
	"goquiz/repository/postgres"
	"log"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbService := os.Getenv("DB_SERVICE")
	dbDSN := os.Getenv("DB_DSN")

	if strings.ToLower(dbService) != DB_Postgres {
		log.Fatalln("only postgres is can be used as scrap datasource")
	}
	model, err := postgres.InitPostgresModel(dbDSN)
	err = postgres.PopulateDB(model)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
