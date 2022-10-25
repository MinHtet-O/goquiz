package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"goquiz/cmd"
	"goquiz/repository/inmemory"
	"goquiz/repository/postgres"
	"goquiz/service"
	"goquiz/transport/rest"
	"log"
)

func main() {
	// load env vars
	err := godotenv.Load()

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
	rateLimiterEnabled, err := cmd.GetenvBool("RATE_LIMITER_ENABLED")
	if err != nil {
		log.Fatalln(err.Error())
	}
	apiTransport, err := cmd.GetenvStr("API_TRANSPORT")
	if err != nil {
		log.Fatalln(err.Error())
	}

	apiKey, err := cmd.GetenvStr("API_KEY")
	if err != nil {
		log.Fatalln(err.Error())
	}
	env, err := cmd.GetenvStr("ENV")
	if err != nil {
		log.Fatalln(err.Error())
	}
	port, err := cmd.GetenvInt("API_PORT")
	if err != nil {
		log.Fatalln(err.Error())
	}
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
	// get configs
	dbService, err := cmd.GetenvStr("DB_SERVICE")
	if err != nil {
		log.Fatalln(err.Error())
	}
	dbDSN, err := cmd.GetenvStr("DB_DSN")
	if err != nil {
		log.Fatalln(err.Error())
	}
	//select model
	var model *service.Model
	if dbService == cmd.DB_Inmemory {
		var err error
		model, err = inmemory.InitInMemoryModel()
		return model, err
	}
	if dbService == cmd.DB_Postgres {
		var err error
		model, err = postgres.InitPostgresModel(dbDSN)
		return model, err
	}
	return nil, fmt.Errorf("unknown model %s", dbService)
}
