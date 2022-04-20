package model

import (
	"time"
)

type Questions []QuestionResp
type Quizzes map[string]Questions

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
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Questions      []QuestionResp `json:"-"`
	QuestionsCount int32          `json:"questions_count,omitempty"`
}

type Question struct {
	ID         int
	WebIndex   int
	Text       string
	AnsOptions []string
	Codeblock  string
	Answer     Answer
	URL        string
}

type QuestionResp struct {
	ID         int      `json:"id"`
	WebIndex   int      `json:"-"`
	Text       string   `json:"text"`
	AnsOptions []string `json:"answers"`
	Codeblock  string   `json:",omitempty"`
	Answer     Answer   `json:"correct_ans"`
	URL        string   `json:"-"`
	Category   Category
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
	question QuestionResp
	ans      Option
	duration time.Time
}
