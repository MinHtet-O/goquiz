package model

import (
	"bufio"
	"fmt"
	"os"
)

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

		fmt.Fprintf(os.Stderr, "Can not create a file for category %s", fileName)
	}

	dataWriter := bufio.NewWriter(file)

	dataWriter.WriteString(content)

	if err := dataWriter.Flush(); err != nil {

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

		return
	}
	if err := createDirectory("questions/" + domain); err != nil {
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
