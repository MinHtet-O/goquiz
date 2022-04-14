package postgres

import (
	"database/sql"
	"fmt"
	"goquiz/pkg/model"
)

type Model struct {
	QuestionsModel
	CategoriesModel
}

type QuestionsModel struct {
	DB *sql.DB
}

type CategoriesModel struct {
	DB *sql.DB
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

		fmt.Println("Successfully Write Category, ", categ, " with ID ", categID)
		fmt.Println("Total QuestionResp ", len(categ.Questions))

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
