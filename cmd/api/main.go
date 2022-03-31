package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"goquiz/pkg/model/postgres"
	"goquiz/pkg/scrapper"
	"log"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"

type config struct {
	scrap bool
	env   string
	port  int
	db    struct {
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
}

type application struct {
	config config
	models postgres.QuizzModel
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	// database related flags
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GOQUIZ_DB_DEV"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	//limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.scrap, "scrap", false, "Scrap the questions and populate db")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Database instance created")
	quizModel := postgres.QuizzModel{DB: db}

	if cfg.scrap { // scrap the questions from the web and populate to database
		s := scrapper.New()
		s.ScrapQuizzes()
		quizModel.InsertQuizzes(s.Quizzes)
	}

	app := &application{
		config: cfg,
		models: quizModel,
	}

	err = app.serve()
	if err != nil {
		log.Fatalln(os.Stdout, err)
	}
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
