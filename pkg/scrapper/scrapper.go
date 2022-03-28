package scrapper

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	model "goquiz/pkg/model"
	"os"
	"regexp"
	"sync"
)

const (
	// TODO: debug why minlimit is not actually limiting the question input
	minlimit = 10 // min question limit threshold for each questions for the category
)

type QuizScrapper struct {
	rootDomain string
	wg         *sync.WaitGroup
	Quizzes    model.Quizzes
	mu         sync.Mutex
	mcqURLs    string
}

// Initialize a new scrapper, currently only javatpoint.com URL is supported
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
	categs := s.getCategories()
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
			s.Quizzes.AddQuestions(categ, *questions)
			s.mu.Unlock()
			// TODO: make save file as dynamic
			//model.SaveQuestionFile(s.rootDomain, categ, fmt.Sprintf("%v", questions))
			fmt.Println("Finish scraping function for ", categ)
		}(c)
	}
	fmt.Println("Waiting to get questions from domain ", s.rootDomain)
	s.wg.Wait()
}

// get categories string arr from mcq URL file
func (s *QuizScrapper) getCategories() []string {
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

func (s *QuizScrapper) scrapQuestions(url string, category string) *[]model.Question {

	var (
		baseUrl = "https://www." + url + "/" + string(category)
		domain  = "www." + url

		questions = make([]model.Question, 0)
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
	)

	c.OnHTML(".pq", func(e *colly.HTMLElement) {
		// initialize the question struct. set the first pq tag as question text
		question := model.Question{
			Text:    parseTitle(e.Text),
			Options: make([]string, model.O_MAX),
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
