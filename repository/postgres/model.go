package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"goquiz/pkg/scraper"
	"goquiz/service"
	"time"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxIdleTime = "15m"
)

type QuestionsModel struct {
	DB *sql.DB
}

type CategoriesModel struct {
	DB *sql.DB
}

func InitPostgresModel(dsn string) (*service.Model, error) {
	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	return &service.Model{
		QuestionsModel:  QuestionsModel{DB: db},
		CategoriesModel: CategoriesModel{DB: db},
	}, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	duration, err := time.ParseDuration(connMaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// populate the repository with scraped questions by categories. It's bulk insert operation.
func InsertQuestionsByCategs(categs []*service.Category, m service.Model) error {
	fmt.Println("## Populating DB ##")
	fmt.Printf("LEN: %d \n", len(categs))
	for _, categ := range categs {
		categID, err := m.CategoriesModel.Insert(*categ)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// later refactor method - InsertQuestions
		for _, question := range categ.Questions {
			_, err := m.QuestionsModel.Insert(categID, question)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
		// TODO: add transaction rollback
	}
	return nil
}

// TODO: separate generic service.model from InsertQuestionsByCategs

// scrap the questions from the web and populate to database
func PopulateDB(model *service.Model) error {
	s := scraper.New()
	categs := s.ScrapQuizzes()
	err := InsertQuestionsByCategs(categs, *model)
	return err
}
