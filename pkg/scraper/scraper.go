package scraper

import (
	"bufio"
	"fmt"
	model "goquiz/pkg/model"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

const (
	// TODO: debug why minlimit is not actually limiting the question input
	minlimit = 10 // min question limit threshold for each questions for the category
)

type QuizScrapper struct {
	rootDomain string
	wg         *sync.WaitGroup
	Quizzes    model.Quizzes
	Categories []model.Category
	mu         sync.Mutex
	mcqURLs    string
}

// Initialize a new scraper, currently only javatpoint.com URL is supported
func New() *QuizScrapper {
	rootDomain := "javatpoint.com"
	mcqURLs := "./files/mcq.txt" // TODO: make mcqURLs dynamic
	return &QuizScrapper{
		rootDomain: rootDomain,
		wg:         &sync.WaitGroup{},
		Quizzes:    model.Quizzes{},
		mu:         sync.Mutex{},
		mcqURLs:    mcqURLs,
	}
}

// TODO: Refactor question methods for domain "javatpoint.com"
// scrap the URLs and get Quizzes for each category
func (s *QuizScrapper) ScrapQuizzes() {
	categs := s.getCategorieTags()
	fmt.Println(categs)
	s.wg.Add(len(categs))
	for _, c := range categs {
		go func(categ string) {
			defer func() {
				s.wg.Done()
				if r := recover(); r != nil {
					fmt.Println("Recover from getting questions, categ", categ)
				}
			}()
			questions := s.scrapQuestions(s.rootDomain, categ)
			if len(*questions) < minlimit {
				fmt.Fprintln(os.Stderr, "Remove ", categ, " from questions with length", len(*questions))
				return
			}
			s.mu.Lock()
			s.AddCategories(categ, *questions)
			s.mu.Unlock()
			// TODO: make save file as dynamic
			//model.SaveQuestionFile(s.rootDomain, categ, fmt.Sprintf("%v", questions))
			fmt.Println("Finish scraping function for ", categ)
		}(c)
	}
	fmt.Println("Waiting to get questions from domain ", s.rootDomain)
	s.wg.Wait()
}

// get categories tag arr from mcq URL file
// tags example - machine-learning-mcq-part1, goang-mcq-part1, dbms-mcq
func (s *QuizScrapper) getCategorieTags() []string {
	categs := []string{}
	re := regexp.MustCompile(`(\w+|\-)+$`)
	file, _ := os.Open(s.mcqURLs)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		categs = append(categs, re.FindString(scanner.Text()))
	}
	return categs
}

func (s *QuizScrapper) scrapQuestions(url string, category string) *[]model.QuestionResp {

	var (
		baseUrl = "https://www." + url + "/" + string(category)
		domain  = "www." + url

		questions = make([]model.QuestionResp, 0)
	)

	defer func() {
		r := recover()
		if r != nil {
			fmt.Fprintf(os.Stderr, "Error fecting %s !", category)
		}
	}()

	// setup new collector
	c := colly.NewCollector(
		colly.AllowedDomains(domain, url),
		colly.Async(true),
	)

	c.OnHTML(".pq", func(e *colly.HTMLElement) {
		// initialize the question struct. set the first pq tag as question text
		question := model.QuestionResp{
			Text:       parseTitle(e.Text),
			AnsOptions: make([]string, model.O_MAX),
		}

		// post with second 'pq' text which need to append into parent question text
		if quesNode := findSibling(e.DOM, "pq", 4, D_Next); quesNode != nil {
			question.Text += "\n" + quesNode.Text()
		}

		// ignore second pq question tags in the post
		if quesNode := findSibling(e.DOM, "pq", 4, D_Prev); quesNode != nil {
			fmt.Fprintf(os.Stderr, "Second pq question tag found category %s , index %d \n", category, e.Index)
			return
		}

		// post that contains codeblock, codeblock is optional
		if codeblockNode := findSibling(e.DOM, "codeblock", 4, D_Next); codeblockNode != nil {
			question.Codeblock = codeblockNode.Text()
		}

		// find answer options for a question
		if err := parseAnsOptions(e, &question, category); err != nil { // no answer options found for this question
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		// correct ans for the question
		if err := parseCorrectAns(e, &question, category); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		question.URL = baseUrl
		// append a new question to question array
		question.WebIndex = e.Index
		questions = append(questions, question)
	})

	c.Visit(baseUrl)
	c.Wait()
	return &questions
}

// add categories together with questions
func (s *QuizScrapper) AddCategories(categ string, ques []model.QuestionResp) {
	// get the category title from input category tag
	// for example get Machine Learning from machine-learning-mcq-part1
	categTitle := func() string {
		tmpArr := strings.Split(categ, "mcq")
		categArr := strings.Split(tmpArr[0], "-")
		for i, val := range categArr {
			categArr[i] = strings.Title(val)
		}
		categTitle := strings.Trim(strings.Join(categArr, " "), " ")
		return categTitle
	}()
	fmt.Println("Category Title", categ)
	// check if there exists questions with the same category title
	found, index := false, 0
	for i, c := range s.Categories {
		if c.Name == categTitle {
			found = true
			index = i
			break
		}
	}
	// combine questions if there are questions with the same categ title
	// for eg - golang-mcq-part1 tag and golang-mcq-part2 tag will be in the same catgory "Golang"

	if found {
		fmt.Printf("Category key %s already exists with len %d for input category %s "+
			"with len %d \n", categTitle, len(categTitle), categ, len(categ))
		s.Categories[index].Questions = append(s.Categories[index].Questions, ques...)
		return
	}

	// Create new category and add questions
	categStruct := model.Category{
		Name:      categTitle,
		Questions: ques,
	}
	s.Categories = append(s.Categories, categStruct)
}

//func (s *QuizScrapper) GetMCQLinks() {
//	baseUrl := "https://" + s.rootDomain + "/"
//	c := colly.NewCollector(
//		colly.AllowedDomains("www.javatpoint.com", "javatpoint.com"),
//	)
//
//	// Find and visit all links
//	c.OnHTML("a", func(e *colly.HTMLElement) {
//		e.Request.Visit(e.Attr("href"))
//	})
//
//	c.OnRequest(func(r *colly.Request) {
//
//		// mcq url found
//		if strings.Contains(r.URL.String(), "mcq") {
//			fmt.Println("FOUND ", r.URL)
//		} else {
//			fmt.Println(r.URL.String())
//		}
//	})
//	c.Visit(baseUrl)
//	c.Wait()
//}
