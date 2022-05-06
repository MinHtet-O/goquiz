package model

import "errors"

func (m *CategoriesModel) GetAll() ([]*Category, error) {

	return *m.Categories, nil
}

func (m *CategoriesModel) GetByID(categId int) (*Category, error) {

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

func (m *CategoriesModel) Insert(categ Category) (int, error) {
	categ.Id = len(*m.Categories) + 1
	*m.Categories = append(*m.Categories, &categ)
	return categ.Id, nil
}
