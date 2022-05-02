package model

import (
	"fmt"
	"time"
)

type Model struct {
	QuestionsModel interface {
		// GetAllByCategoryId(categId int) ([]m.Question, error)
		GetAll(category Category) ([]Question, error)
		Insert(categID int, q Question) error
	}
	CategoriesModel interface {
		GetAll() ([]*Category, error)
		GetByID(categId int) (*Category, error)
		Insert(cate Category) (int, error)
	}
}

type QuestionsModel struct{ Categories []*Category }
type CategoriesModel struct{ Categories []*Category }

type Questions []Question

type Option int8

const (
	A Option = iota
	B
	C
	D
	E
	O_MAX
)

var AnsMapping = map[string]Option{
	"a": A,
	"b": B,
	"c": C,
	"d": D,
	"e": E,
}

type Category struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Questions      []Question `json:"-"`
	QuestionsCount int32      `json:"questions_count,omitempty"`
}

//type Question struct {
//	ID         int
//	WebIndex   int
//	Text       string
//	AnsOptions []string
//	Codeblock  string
//	Answer     Answer
//	URL        string
//}

type Question struct {
	ID         int      `json:"id"`
	WebIndex   int      `json:"-"`
	Text       string   `json:"text"`
	AnsOptions []string `json:"answers"`
	Codeblock  string   `json:",omitempty"`
	Answer     Answer   `json:"correct_ans"`
	URL        string   `json:"-"`
	// TODO: remove Category from Question
	Category Category `json:"-"`
}

type Answer struct {
	Option      Option `json:"option"`
	Explanation string `json:"explanation"`
}

//Optional data structure
type User struct {
	name string
}

// add method to change the setting - setters
type Setting struct {
	quesTimeout int
}

//MCQ run time data structure
// add method to calculate totalScore
type MatchRecord struct {
	choices    []Choice
	date       time.Time
	totalScore int
}

type Choice struct {
	question Question
	ans      Option
	duration time.Time
}

func (m Model) InsertCategories(categs []*Category) error {
	fmt.Println("Insert CategoriesModel")
	fmt.Printf("LEN: %d \n", len(categs))
	for _, categ := range categs {

		categID, err := m.CategoriesModel.Insert(*categ)
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
