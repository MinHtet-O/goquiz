package model

import (
	"time"
)

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
