package main

import (
	"context"
	"database/sql"
	"flag"
	m "goquiz/pkg/model"
	"goquiz/pkg/model/postgres"
	"goquiz/pkg/scraper"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	populateDB bool
	env        string
	port       int
	db         struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	// rate limiting for each ip address
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}

	auth struct {
		enabled bool
		apiKey  string
	}
}

type application struct {
	config config
	models m.Model
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	// database related params
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv(""), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	//rate limiter params
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	//other params
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.populateDB, "populate-db", false, "Populate the database")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	//authenticaiton
	flag.StringVar(&cfg.auth.apiKey, "apikey", "", "API-key for Authentication")
	flag.Parse()

	var model m.Model
	model, err := setModel(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if cfg.populateDB { // scrap the questions from the web and populate to database

		if cfg.db.dsn == "" {
			log.Fatal("please provide database service name to populate")
		}
		s := scraper.New()
		s.ScrapQuizzes()
		// TODO: type cast to postgres.Model
		model.InsertCategories(s.Categories)
		os.Exit(1)
	}

	app := &application{
		config: cfg,
		models: model,
	}

	err = app.serve()
	if err != nil {
		log.Fatalln(os.Stdout, err)
	}
}

func setModel(cfg config) (m.Model, error) {

	if cfg.db.dsn == "" {
		// scrap the questions in db-less mode
		s := scraper.New()
		s.ScrapQuizzes()

		// return in memory model, only pointers of in memory model implement model interfaces.
		// methods need pointer receiver as they need to edit in memory data in place
		return m.Model{
			QuestionsModel:  &m.QuestionsModel{&s.Categories},
			CategoriesModel: &m.CategoriesModel{&s.Categories},
		}, nil
	}

	db, err := openDB(cfg)
	if err != nil {
		return m.Model{}, err
	}
	// return postgres model
	return m.Model{
		QuestionsModel:  postgres.QuestionsModel{DB: db},
		CategoriesModel: postgres.CategoriesModel{DB: db},
	}, nil
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
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

// TODO: figure out why scrap results are different each time the program run
// TODO: change standard logging to json format, replace all log stdout stderr with json logging
// TODO: why file write is not working
// TODO: move URL field from questions to categories
// TODO: make category table composite primary key (id,name). So, there would be no duplicate category name
