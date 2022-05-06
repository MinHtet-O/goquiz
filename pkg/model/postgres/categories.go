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

func (m CategoriesModel) GetByID(categId int) (*model.Category, error) {
	query := `select id,name from categories where id=$1`
	var category model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, categId).Scan(&category.Id, &category.Name)

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
	fmt.Println("Category Id")
	fmt.Println(categ.id)
	if err != nil {
		return 0, err
	}
	return categ.id, nil
}

func (m Model) InsertCategories(categs []model.Category) error {
	fmt.Println("Insert CategoriesModel")
	fmt.Printf("LEN: %d \n", len(categs))
	for _, categ := range categs {

		categID, err := m.CategoriesModel.Insert(categ)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// later refactor method - InsertQuestions
		for _, question := range categ.Questions {
			err := m.QuestionsModel.Insert(categID, question)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
		// TODO: add transaction rollback
		//if err != nil {
		//	// rollback transaction
		//	continue
		//}
		// commit transaction
	}
	return nil
}
