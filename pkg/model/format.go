package model

import (
	"fmt"
	"strings"
)

func (q Quizzes) String() string {
	var sb strings.Builder
	for key, val := range q {
		fmt.Fprintf(&sb, "key %s ,value %v \n", key, len(val))
	}
	return sb.String()
}

func (questions Questions) String() string {
	var sb strings.Builder
	for i, q := range questions {
		fmt.Fprintln(&sb, "No.", i)
		fmt.Fprintln(&sb, q)
	}
	return sb.String()
}

func (a Answer) String() string {
	return fmt.Sprintf("Option ", a.Option, " ", a.Explanation)
}

// TODO: Why to string not working with *QuestionResp ??
func (q QuestionResp) String() string {
	var opsb strings.Builder
	for i, val := range q.AnsOptions {
		fmt.Fprintln(&opsb, i, ".", val)
	}
	return fmt.Sprintf("Index: %d\nQuestion: %s\nCodeBlock:\n%s\nOptions:\n%vCorrect Answer: %d\n\n%s\n\n", q.WebIndex, q.Text, q.Codeblock, opsb.String(), q.Answer.Option, q.Answer.Explanation)
}
