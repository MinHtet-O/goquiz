package main

import (
	"fmt"
	"goquiz/repository/inmemory"
	"goquiz/repository/postgres"
	"goquiz/service"
	"goquiz/transport/rest"
	"log"
)

func main() {
	// parse the cli configs
	var cfg Config
	parseConfig(&cfg)

	// select model
	model, err := setModel(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// select transport
	app, err := setTransport(cfg, model)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = app.Serve()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func setTransport(cfg Config, model *service.Model) (service.Transport, error) {
	// select transport protocol
	if cfg.Transport == Transport_REST {
		return rest.NewRESTServer(model), nil
	}

	if cfg.Transport == Transport_GRPC {
		return nil, fmt.Errorf("GRPC Transport is still in development")
	}
	return nil, fmt.Errorf("Unknown Transport type %s", cfg.Transport)
}

func setModel(cfg Config) (*service.Model, error) {
	var model *service.Model
	if cfg.Db.Name == DB_Postgres {
		var err error
		model, err = inmemory.InitInMemoryModel()
		return model, err
	}
	if cfg.Db.Name == DB_Inmemory {
		var err error
		model, err = postgres.InitPostgresModel(
			cfg.Db.Dsn,
			cfg.Db.MaxOpenConns,
			cfg.Db.MaxIdleConns,
			cfg.Db.MaxIdleTime,
		)
		return model, err
	}
	return nil, fmt.Errorf("unknown model %s", cfg.Db.Name)
}
