package model

import (
	"fmt"
	"strings"
)

//MCQ based data structure

type Questions []Question
type Quizzes map[string]Questions

func (q Quizzes) String() string {
	var sb strings.Builder
	for key, val := range q {
		fmt.Fprintf(&sb, "key %s ,value %v \n", key, len(val))
	}
	return sb.String()
}

func (q Quizzes) AddQuestions(categ string, ques Questions) {

	// get the category title from input category
	categTitle := func() string {
		tmpArr := strings.Split(categ, "mcq")
		categArr := strings.Split(tmpArr[0], "-")
		for i, val := range categArr {
			categArr[i] = strings.Title(val)
		}
		categTitle := strings.Trim(strings.Join(categArr, " "), " ")
		return categTitle
	}()

	// combine questions if there are questions with the same categ title
	if val, ok := q[categTitle]; ok {
		fmt.Printf("Category %s alredy exists with len %d for input category %s with len %d \n", categTitle, len(categTitle), categ, len(categ))
		val = append(val, ques...)
		q[categTitle] = val
	} else {
		q[categTitle] = ques
	}
}

func (questions Questions) String() string {
	var sb strings.Builder
	for i, q := range questions {
		fmt.Fprintln(&sb, "No.", i)
		fmt.Fprintln(&sb, q)
	}
	return sb.String()
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

func (a Answer) String() string {
	return fmt.Sprintf("Option ", a.Option, " ", a.Explanation)
}

// TODO: Why to string not working with *Question ??
func (q Question) String() string {
	var opsb strings.Builder
	for i, val := range q.Options {
		fmt.Fprintln(&opsb, i, ".", val)
	}
	return fmt.Sprintf("Index: %d\nQuestion: %s\nCodeBlock:\n%s\nOptions:\n%vCorrect Answer: %d\n\n%s\n\n", q.WebIndex, q.Text, q.Codeblock, opsb.String(), q.CorrectAns.Option, q.CorrectAns.Explanation)
}
