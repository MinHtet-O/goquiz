package utils

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"time"
)

type nodeDirection int8

const (
	// Set the direction to find from the current node
	D_Next nodeDirection = iota // find in the next nodes from the current node
	D_Prev                      // find in the previous nodes from the current node
)

func FindSibling(node *goquery.Selection, class string, depth int, dir nodeDirection) *goquery.Selection {

	for i := 1; i < depth; i++ {
		if dir == D_Next {
			node = node.Next()
		} else {
			node = node.Prev()
		}
		nodeclass, exist := node.Attr("class")
		if exist && nodeclass == class {
			return node
		}
	}
	return nil

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

func ValidateImgURL(url string) error {
	client := clientHTTP()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	return nil
}

// TODO: set one http client instance for each scrapper
func clientHTTP() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: t,
	}
	return client
}
