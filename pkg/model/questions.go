package model

import (
	"errors"
)

func (m QuestionsModel) GetAll(category Category) ([]Question, error) {

	found, foundIndex := false, 0

	for index, categ := range m.Categories {
		if category.ID == categ.ID {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		return []Question{}, errors.New("no record found")
	}

	return m.Categories[foundIndex].Questions, nil
}

func (m QuestionsModel) Insert(categID int, q Question) error {
	return nil
}
