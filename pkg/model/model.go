package model

import (
	"time"
)

type Questions []Question
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

type Question struct {
	WebIndex   int
	Text       string
	Options    []string
	Codeblock  string
	CorrectAns Answer
}

type Answer struct {
	Option      Option
	Explanation string
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
