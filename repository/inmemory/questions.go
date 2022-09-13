package inmemory

import (
	"errors"
	"goquiz/service"
)

func (m *QuestionsModel) GetAll(category service.Category) ([]service.Question, error) {

	found, foundIndex := false, 0

	for index, categ := range *m.Categories {
		if category.Id == categ.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		return []service.Question{}, errors.New("no record found")
	}

	return (*m.Categories)[foundIndex].Questions, nil
}

func (m *QuestionsModel) Insert(categID int, que service.Question) (int, error) {

	found, foundIndex := false, 0

	for index, categ := range *m.Categories {
		if categID == categ.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		return 0, errors.New("no record found")
	}
	que.Id = len((*m.Categories)[foundIndex].Questions) + 1
	(*m.Categories)[foundIndex].Questions = append((*m.Categories)[foundIndex].Questions, que)

	return que.Id, nil
}
