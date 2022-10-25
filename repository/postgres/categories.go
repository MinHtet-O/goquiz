package postgres

import (
	"context"
	"database/sql"
	"errors"
	"goquiz/service"
	"time"
)

func (m CategoriesModel) GetAll() ([]*service.Category, error) {
	query := `select c.id, c."name", (select count(*) from questions WHERE questions.category_id=c.id) FROM categories c`
	var categories []*service.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category service.Category
		err := rows.Scan(&category.Id, &category.Name, &category.QuestionsCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	// TODO: why do we need to close rows??
	defer rows.Close()

	return categories, nil
}

func (m CategoriesModel) GetByID(categId int) (*service.Category, error) {
	query := `select id,name from categories where id=$1`
	var category service.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, categId).Scan(&category.Id, &category.Name)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, service.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &category, nil
}

func (m CategoriesModel) Insert(cate service.Category) (int, error) {
	query := `INSERT INTO categories (name) values ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categ := struct {
		id int
	}{}

	err := m.DB.QueryRowContext(ctx, query, cate.Name).Scan(&categ.id)

	if err != nil {
		return 0, err
	}
	return categ.id, nil
}
