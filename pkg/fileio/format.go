package fileio

import (
	"fmt"
	"goquiz/service"
	"strings"
)

func PrintQuestions(questions service.Questions) string {
	var sb strings.Builder
	for i, q := range questions {
		fmt.Fprintln(&sb, "No.", i)
		fmt.Fprintln(&sb, q)
	}
	return sb.String()
}

func PrintAnswer(a service.Answer) string {
	return fmt.Sprintf("Option ", a.Option, " ", a.Explanation)
}

func PrintQuestion(q service.Question) string {
	var opsb strings.Builder
	for i, val := range q.AnsOptions {
		fmt.Fprintln(&opsb, i, ".", val)
	}
	return fmt.Sprintf("Index: %d\nQuestion: %s\nCodeBlock:\n%s\nOptions:\n%vCorrect Answer: %d\n\n%s\n\n", q.WebIndex, q.Text, q.Codeblock, opsb.String(), q.Answer.Option, q.Answer.Explanation)
}
