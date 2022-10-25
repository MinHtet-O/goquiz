package main

import (
	"github.com/joho/godotenv"
	"goquiz/cmd"
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

	dbService := cmd.GetenvStr("DB_SERVICE", "postgres")
	dbDSN := cmd.GetenvStr("DB_DSN", "")
	if dbDSN == "" {
		log.Fatalln("db dsn couldn't be empty")
		os.Exit(1)
	}
	if strings.ToLower(dbService) != cmd.DB_Postgres {
		log.Fatalln("only postgres is can be used as scrap datasource")
	}
	model, err := postgres.InitPostgresModel(dbDSN)
	err = postgres.PopulateDB(model)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
