package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"goquiz/pkg/scraper"
	"goquiz/repository/inmemory"
	"goquiz/repository/postgres"
	"goquiz/service"
	"goquiz/transport/rest"
	"log"
	"os"
	"time"
)

func main() {
	var cfg service.Config
	parseConfig(&cfg)

	// select model/ repository
	model, err := SetModel(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// select transport
	app, err := SetTransport(cfg, model)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = app.Serve()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//if cfg.PopulateDB { // scrap the questions from the web and populate to database
	//
	//	if cfg.Db.Dsn == "" {
	//		log.Fatal("please provide database service name to populate")
	//	}
	//	s := scraper.New()
	//	s.ScrapQuizzes()
	//	// TODO: type cast to postgres.PostgresModel
	//	err = model.InsertQuestionsByCategs(s.Categories)
	//	if err != nil {
	//		log.Fatalln(err.Error())
	//	}
	//	os.Exit(1)
	//}

}

func parseConfig(cfg *service.Config) {
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
	flag.StringVar(&cfg.Transport, "transport", service.Transport_REST, "Transport protocol for your API ( rest | grpc )")
	flag.Parse()
}

func SetTransport(cfg service.Config, model service.Model) (service.Transport, error) {
	// select transport protocol
	if cfg.Transport == service.Transport_REST {
		return rest.NewRESTServer(cfg, model), nil
	}

	if cfg.Transport == service.Transport_GRPC {
		return nil, fmt.Errorf("GRPC Transport is still in development")

	}
	return nil, fmt.Errorf("Unknown Transport type %s", cfg.Transport)
}

func SetModel(cfg service.Config) (service.Model, error) {

	if cfg.Db.Dsn == "" {
		// scrap the questions in db-less mode
		s := scraper.New()
		s.ScrapQuizzes()

		// return in inmemory model, only pointers of in inmemory model implement model interfaces.
		// methods need pointer receiver as they need to edit in inmemory data in place
		return service.Model{
			QuestionsModel:  &inmemory.QuestionsModel{&s.Categories},
			CategoriesModel: &inmemory.CategoriesModel{&s.Categories},
		}, nil
	}

	db, err := OpenDB(cfg)
	if err != nil {
		return service.Model{}, err
	}
	// return postgres model
	return service.Model{
		QuestionsModel:  postgres.QuestionsModel{DB: db},
		CategoriesModel: postgres.CategoriesModel{DB: db},
	}, nil
}

func OpenDB(cfg service.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
