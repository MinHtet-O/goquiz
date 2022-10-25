package scrap

import (
	"fmt"
	"goquiz/repository/inmemory"
	"goquiz/repository/postgres"
	"goquiz/service"
	"log"
)

func test() {
	// parse the cli configs
	var cfg Config
	parseConfig(&cfg)

	// select model
	model, err := setModel(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = postgres.PopulateDB(model)
	if err != nil {
		log.Fatalln(err.Error())
	}
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
