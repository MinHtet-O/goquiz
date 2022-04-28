package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"goquiz/pkg/model"
	"regexp"
	"strings"
	"unicode"
)

func parseCorrectAns(e *colly.HTMLElement, question *model.QuestionResp, category string) error {

	ansReg := regexp.MustCompile(`(\(\w\))|(\s\w\.$)|(\s\w\)\s)|(\s\w$)|(\s\w\.\s\w*)|(\s\w:\s)|(\s\w\s)`)
	// find the answer node
	ansNode := findSibling(e.DOM, "testanswer", 8, D_Next)

	if ansNode == nil {
		return fmt.Errorf("No answer node found %s , index %d \n", category, e.Index)
	}

	// Index 0, first child is the correct answer
	ansTxt := ansNode.Children().First().Text()
	ans := ansReg.FindString(ansTxt)

	if ans == "" {
		return fmt.Errorf("Invalid Correct Answer Format %s , index %d, %s \n", category, e.Index, ansTxt)
	}

	// make necessary string processing to extract the answer option
	ans = strings.ToLower(strings.Split(ans, "")[1])
	if opt, found := model.AnsMapping[ans]; found {
		question.Answer.Option = opt
	}

	ansNode.Children().Each(func(i int, c *goquery.Selection) {
		if i != 0 {
			// index 1 to ... is explanation
			// append the explanation
			question.Answer.Explanation += c.Text()
		}
	})

	return nil
}

func parseAnsOptions(e *colly.HTMLElement, question *model.QuestionResp, category string) error {

	// answer options as text
	// valid A,B,C,D,E options
	optsNode := findSibling(e.DOM, "pointsa", 6, D_Next)
	if optsNode != nil && optsNode.Children().Length() < 6 {
		optsNode.Children().Each(func(i int, c *goquery.Selection) {
			question.AnsOptions[model.Option(i)] = strings.ToLower(c.Text())
		})
		return nil
	}

	// answer options as image with URL
	imageURL, exist := e.DOM.Next().Children().First().Attr("src")
	if exist {
		// test the image URL to make sure it is the valid URL
		err := validateImageURL(imageURL)
		if err != nil {
			return fmt.Errorf("Invalid Image URL %s, for category %s \n", imageURL, category)
		}
		question.AnsOptions[0] = imageURL
		return nil
	}

	// no anser option found for the given question
	return fmt.Errorf("No Answer Option found for category %s , index %d \n", category, e.Index)
}

func parseTitle(title string) string {
	// check the string the see if it starts with num)
	valid := true
	breakIndex := 0
	for i, r := range title {
		if unicode.IsDigit(r) || r == ' ' {
			continue
		}
		if r == ')' {
			breakIndex = i
			break
		}
		valid = false
		break
	}
	if valid {
		title = title[breakIndex+2:]
	}
	return title
}
