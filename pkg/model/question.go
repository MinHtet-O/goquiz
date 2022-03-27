package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (q Quizzes) String() string {
	var sb strings.Builder
	for key, val := range q {
		fmt.Fprintf(&sb, "key %s ,value %v \n", key, len(val))
	}
	return sb.String()
}

func (q Quizzes) AddQuestions(categ string, ques Questions) {

	// get the category title from input category URL
	categTitle := func() string {
		tmpArr := strings.Split(categ, "mcq")
		categArr := strings.Split(tmpArr[0], "-")
		for i, val := range categArr {
			categArr[i] = strings.Title(val)
		}
		categTitle := strings.Trim(strings.Join(categArr, " "), " ")
		return categTitle
	}()

	// combine questions if there are questions with the same categ key
	if val, ok := q[categTitle]; ok {
		fmt.Printf("Category %s alredy exists with len %d for input category %s with len %d \n", categTitle, len(categTitle), categ, len(categ))
		val = append(val, ques...)
		q[categTitle] = val
		return
	}

	q[categTitle] = ques

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

// TODO: Why to string not working with *Question ??
func (q Question) String() string {
	var opsb strings.Builder
	for i, val := range q.Options {
		fmt.Fprintln(&opsb, i, ".", val)
	}
	return fmt.Sprintf("Index: %d\nQuestion: %s\nCodeBlock:\n%s\nOptions:\n%vCorrect Answer: %d\n\n%s\n\n", q.WebIndex, q.Text, q.Codeblock, opsb.String(), q.CorrectAns.Option, q.CorrectAns.Explanation)
}

func SaveText(fileName string, content string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "Recover from file write failure for catgory", fileName)
		}
	}()

	//delete remaining file
	if err := os.Truncate(fileName, 0); err != nil {
		//fmt.Fprintln(os.Stderr, "Failed to truncate the file for category %s : %v",
		//	fileName, err)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(os.Stderr, "Can not create a file for category %s", fileName)
	}

	dataWriter := bufio.NewWriter(file)

	dataWriter.WriteString(content)

	if err := dataWriter.Flush(); err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(os.Stdout, "Can not create a file for category %s", fileName)
	}
}

func SaveQuestionToFile(domain string, category string, content string) {
	// TODO: refactor question const variable
	path := fmt.Sprintf("questions/%s", domain)
	file := fmt.Sprintf("%s/%s.txt", path, category)

	if err := creteDirectory("questions"); err != nil {
		fmt.Println(domain, err.Error())
		return
	}
	if err := creteDirectory("questions/" + domain); err != nil {
		fmt.Println(domain+"/"+path, err.Error())
		return
	}
	SaveText(file, content)
}

func creteDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
