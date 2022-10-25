package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"goquiz/cmd"
	"goquiz/repository/inmemory"
	"goquiz/repository/postgres"
	"goquiz/service"
	"goquiz/transport/rest"
	"log"
	"strings"
)

func main() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// select model
	model, err := setModel()
	if err != nil {
		log.Fatalln(err.Error())
	}
	// select transport
	app, err := setTransport(model)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// run the app
	err = app.Serve()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func setTransport(model *service.Model) (service.Transport, error) {

	// get configs
	rateLimiterEnabled := cmd.GetenvBool("RATE_LIMITER_ENABLED", true)
	apiTransport := cmd.GetenvStr("API_TRANSPORT", "rest")
	apiKey := cmd.GetenvStr("API_KEY", "")
	env := cmd.GetenvStr("ENV", "dev")
	port := cmd.GetenvInt("API_PORT", 8080)

	// select transport protocol
	if apiTransport == cmd.Transport_REST {
		config := rest.Config{
			apiKey,
			rateLimiterEnabled,
			env,
			port,
		}
		transport := rest.NewRESTServer(model, config, env)
		return transport, nil
	}

	if apiTransport == cmd.Transport_GRPC {
		return nil, fmt.Errorf("GRPC Transport is still in development")
	}
	return nil, fmt.Errorf("Unknown Transport type %s", apiTransport)
}

func setModel() (*service.Model, error) {
	var model *service.Model
	// get configs
	dbService := cmd.GetenvStr("DB_SERVICE", "postgres")
	dbDSN := cmd.GetenvStr("DB_DSN", "")
	if dbDSN == "" {
		return model, errors.New("db dsn couldn't be empty")
	}
	//select model
	if strings.ToLower(dbService) == cmd.DB_Inmemory {
		var err error
		model, err = inmemory.InitInMemoryModel()
		return model, err
	}
	if strings.ToLower(dbService) == cmd.DB_Postgres {
		var err error
		model, err = postgres.InitPostgresModel(dbDSN)
		return model, err
	}
	return nil, fmt.Errorf("unknown model %s", dbService)
}
