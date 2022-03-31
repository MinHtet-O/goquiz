package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	// for eg - golang-mcq-part1 and golang-mcq-part2 will be in the same map entry
	if val, ok := q[categTitle]; ok {
		fmt.Printf("Category key %s already exists with len %d for input category %s with len %d \n", categTitle, len(categTitle), categ, len(categ))
		val = append(val, ques...)
		q[categTitle] = val
		return
	}

	q[categTitle] = ques
}

func SaveFile(fileName string, content string) {
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

//TODO: method not working
//TODO: save as csv format
func SaveQuestionFile(domain string, category string, content string) {
	// TODO: refactor question const variable
	path := fmt.Sprintf("files/questions/%s", domain)
	file := fmt.Sprintf("%s/%s.txt", path, category)

	if err := createDirectory("questions"); err != nil {
		fmt.Println(domain, err.Error())
		return
	}
	if err := createDirectory("questions/" + domain); err != nil {
		fmt.Println(domain+"/"+path, err.Error())
		return
	}
	SaveFile(file, content)
}

func createDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
