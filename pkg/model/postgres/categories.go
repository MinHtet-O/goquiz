package postgres

import (
	"context"
	"goquiz/pkg/model"
	"time"
)

func (m CategoriesModel) GetCategories() ([]*model.Category, error) {
	query := `select c.id, c."name", (select count(*) from questions WHERE questions.category_id=c.id) FROM categories c`
	var categories []*model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name, &category.QuestionsCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	// TODO: why do we need to close rows??
	defer rows.Close()
	defer cancel()
	return categories, nil
}
