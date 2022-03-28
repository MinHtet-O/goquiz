package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"goquiz/pkg/model/postgres"
	"goquiz/pkg/scrapper"
	"log"
	"os"
	"time"
)

type config struct {
	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func main() {
	var cfg config

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GOQUIZ_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	db, err := openDB(cfg)

	if err != nil {
		log.Fatalln(err.Error())
	}

	s := scrapper.New()
	s.ScrapQuizzes()
	//s.GetMCQLinks()

	quizModel := postgres.QuizzModel{DB: db}
	quizModel.InsertQuizzes(s.Quizzes)
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
	// Use PingContext() to establish a new connection to the database,
	//If the connection couldn't be established successfully within the 5 second deadline,
	//then this will return an  error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
