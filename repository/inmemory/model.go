package inmemory

import (
	"goquiz/service"
)

type MemoryModel struct {
	QuestionsModel
	CategoriesModel
}
type QuestionsModel struct{ Categories *[]*service.Category }
type CategoriesModel struct{ Categories *[]*service.Category }
