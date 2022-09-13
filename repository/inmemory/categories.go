package inmemory

import (
	"errors"
	"goquiz/service"
)

func (m *CategoriesModel) GetAll() ([]*service.Category, error) {

	return *m.Categories, nil
}

func (m *CategoriesModel) GetByID(categId int) (*service.Category, error) {

	found, foundIndex := false, 0
	for index, categ := range *m.Categories {
		if categId == categ.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		return nil, errors.New("no record found")
	}

	return (*m.Categories)[foundIndex], nil
}

func (m *CategoriesModel) Insert(categ service.Category) (int, error) {
	categ.Id = len(*m.Categories) + 1
	*m.Categories = append(*m.Categories, &categ)
	return categ.Id, nil
}
