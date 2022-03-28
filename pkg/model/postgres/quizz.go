package postgres

import (
	"context"
	"fmt"
	"goquiz/pkg/model"
	"time"
)

func (m QuizzModel) InsertQuizzes(q model.Quizzes) error {
	fmt.Println("LENGTH!!")
	fmt.Println(len(q))
	i := 0
	for categ, _ := range q {
		i++
		categID, err := m.InsertCategory(categ)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Print("Successfully Write Category, ", categ, " with ID", categID)
		//for _, question := range questions {
		//	m.InsertQuestion(categ, question)
		//}
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

	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (m QuizzModel) InsertQuestion(category string, q model.Question) error {
	return nil
}
