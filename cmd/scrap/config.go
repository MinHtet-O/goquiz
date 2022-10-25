package main

import (
	"flag"
	"os"
)

const (
	Transport_REST string = "rest"
	Transport_GRPC string = "grpc"
	DB_Postgres    string = "postgres"
	DB_Inmemory    string = "inmemory"
)

type Config struct {
	PopulateDB bool
	Env        string
	Port       int
	Db         struct {
		Name         string
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	// rate limiting for each ip address
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}

	Auth struct {
		Enabled bool
		ApiKey  string
	}
	Transport string
}

func parseConfig(cfg *Config) {
	// database related params
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv(""), "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	//rate limiter params
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	//other params
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.PopulateDB, "populate-db", false, "Populate the database")
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")

	//authenticaiton
	flag.StringVar(&cfg.Auth.ApiKey, "apikey", "", "API-key for Authentication")

	//transport
	flag.StringVar(&cfg.Transport, "transport", Transport_REST, "Transport protocol for your API ( rest | grpc )")
	flag.Parse()
}
