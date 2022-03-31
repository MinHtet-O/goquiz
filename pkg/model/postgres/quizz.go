package postgres

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"goquiz/pkg/model"
	"time"
)

func (m QuizzModel) InsertQuizzes(quizzes model.Quizzes) error {
	i := 0
	for categ, questions := range quizzes {
		i++
		categID, err := m.InsertCategory(categ)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("Successfully Write Category, ", categ, " with ID ", categID)
		fmt.Println("Total Question ", len(questions))

		for _, question := range questions {
			err := m.InsertQuestion(categID, question)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
		//if err != nil {
		//	// rollback transaction
		//	continue
		//}
		// commit transaction
	}
	fmt.Println("I!!")
	fmt.Println(i)
	return nil
}

func (m QuizzModel) InsertCategory(category string) (int, error) {
	query := `INSERT INTO categories (name) values ($1) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categ := struct {
		id int
	}{}

	err := m.DB.QueryRowContext(ctx, query, category).Scan(&categ.id)
	fmt.Println("Category ID")
	fmt.Println(categ.id)
	if err != nil {
		return 0, err
	}
	return categ.id, nil
}

func (m QuizzModel) InsertQuestion(categID int, q model.Question) error {

	args := []interface{}{categID, q.WebIndex, q.Text, pq.Array(q.Options), q.Codeblock, q.CorrectAns.Option, q.CorrectAns.Explanation, q.URL}
	// TODO: replace query with sql prepared statement
	query := `INSERT INTO questions (category_id,web_index,text,ans_options,code_block,correct_ans_opt,correct_ans_explanation,url) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}
	return nil
}

func (m QuizzModel) GetCategories() ([]*model.Category, error) {
	query := `SELECT name FROM categories`
	var categories []*model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Name)
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
