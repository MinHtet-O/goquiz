package inmemory

import (
	"goquiz/pkg/scraper"
	"goquiz/service"
)

type QuestionsModel struct{ Categories *[]*service.Category }
type CategoriesModel struct{ Categories *[]*service.Category }

func InitInMemoryModel() (*service.Model, error) {
	// scrap the questions
	s := scraper.New()
	categs := s.ScrapQuizzes()
	return &service.Model{
		QuestionsModel:  &QuestionsModel{&categs},
		CategoriesModel: &CategoriesModel{&categs},
	}, nil
}
