package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goquiz/pkg/model"
	"time"
)

func (m CategoriesModel) GetAll() ([]*model.Category, error) {
	query := `select c.id, c."name", (select count(*) from questions WHERE questions.category_id=c.id) FROM categories c`
	var categories []*model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
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

	return categories, nil
}

func (m CategoriesModel) Get(categId int) (*model.Category, error) {
	query := `select id,name from categories where id=$1`
	var category model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, categId).Scan(&category.ID, &category.Name)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoriesModel) Insert(cate model.Category) (int, error) {
	query := `INSERT INTO categories (name) values ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categ := struct {
		id int
	}{}

	err := m.DB.QueryRowContext(ctx, query, cate.Name).Scan(&categ.id)
	fmt.Println("Category ID")
	fmt.Println(categ.id)
	if err != nil {
		return 0, err
	}
	return categ.id, nil
}
