package scrapper

// TODO: refactor for proper naming and package
import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	model "goquiz/pkg/model"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

const (
	// TODO: debug why minlimit is not actually limiting the question input
	minlimit = 10 // min question limit threshold for each questions for the category
)

type QuizScrapper struct {
	rootDomain string
	wg         *sync.WaitGroup
	quizzes    model.Quizzes
	mu         sync.Mutex
	filepath   string
}

// Initialize a new scrapper, currently only javatpoint.com URL is supported
func New() *QuizScrapper {
	rootDomain := "javatpoint.com"
	filepath := "./mcq_url.txt" // TODO: make filepath dynamic
	return &QuizScrapper{
		rootDomain: rootDomain,
		wg:         &sync.WaitGroup{},
		quizzes:    model.Quizzes{},
		mu:         sync.Mutex{},
		filepath:   filepath,
	}
}

// TODO: Refactor question methods for domain "javatpoint.com"
func (s *QuizScrapper) GetQuizzes() {
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
			questions := s.getQuestions(s.rootDomain, categ)
			if len(*questions) < minlimit {
				fmt.Fprintln(os.Stderr, "Remove ", categ, " from questions with length", len(*questions))
				return
			}
			s.mu.Lock()
			s.quizzes.AddQuestions(categ, *questions)
			s.mu.Unlock()
			// TODO: make save file as dynamic
			//model.SaveQuestionToFile(s.rootDomain, categ, fmt.Sprintf("%v", questions))
			fmt.Println("Finish scraping function for ", categ)
		}(c)
	}
	fmt.Println("Waiting to get questions from domain ", s.rootDomain)
	s.wg.Wait()
}

func (s *QuizScrapper) getCategories() []string {
	categs := []string{}
	re := regexp.MustCompile(`(\w+|\-)+$`)
	file, _ := os.Open(s.filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		categs = append(categs, re.FindString(scanner.Text()))
	}
	return categs
}

func (s *QuizScrapper) getQuestions(url string, category string) *[]model.Question {

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
		colly.Async(true),
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

		// append a new question to question array
		question.WebIndex = e.Index
		questions = append(questions, question)
	})

	c.Visit(baseUrl)
	c.Wait()
	return &questions
}

func parseCorrectAns(e *colly.HTMLElement, question *model.Question, category string) error {

	ansReg := regexp.MustCompile(`(\(\w\))|(\s\w\.$)|(\s\w\)\s)|(\s\w$)|(\s\w\.\s\w*)|(\s\w:\s)|(\s\w\s)`)
	// find the answer node
	ansNode := findSibling(e.DOM, "testanswer", 8, D_Next)

	if ansNode == nil {
		return fmt.Errorf("Invalid Correct Answer ", category, ", index", e.Index)
	}

	// Index 0, first child is the correct answer
	ans := ansNode.Children().First().Text()
	ans = ansReg.FindString(ans)

	if ans == "" {
		return fmt.Errorf("Invalid Correct Answer Format ", category, ", index", e.Index)
	}

	// make necessary string processing to extract the answer option
	ans = strings.ToLower(strings.Split(ans, "")[1])
	if opt, found := model.AnsMapping[ans]; found {
		question.CorrectAns.Option = opt
	}

	ansNode.Children().Each(func(i int, c *goquery.Selection) {
		if i != 0 {
			// index 1 to ... is explanation
			// append the explanation
			question.CorrectAns.Explanation += c.Text()
		}
	})

	return nil
}

func parseAnsOptions(e *colly.HTMLElement, question *model.Question, category string) error {

	// answer options as text
	// valid A,B,C,D,E options
	optsNode := findSibling(e.DOM, "pointsa", 6, D_Next)
	if optsNode != nil && optsNode.Children().Length() < 6 {
		optsNode.Children().Each(func(i int, c *goquery.Selection) {
			question.Options[model.Option(i)] = strings.ToLower(c.Text())
		})
		return nil
	}

	// answer options as image with URL
	imageURL, exist := e.DOM.Next().Children().First().Attr("src")
	if exist {
		// test the image URL to make sure it is the valid URL
		err := validateImageURL(imageURL)
		if err != nil {
			return fmt.Errorf("Invalid Image URL", imageURL, " for category ", category)
		}
		question.Options[0] = imageURL
		return nil
	}

	// no anser option found for the given question
	return fmt.Errorf("No Answer Option found for ", category, ", index", e.Index)
}

func (s *QuizScrapper) getMCQLinks() {
	baseUrl := "https://www.javatpoint.com/"
	c := colly.NewCollector(
		colly.AllowedDomains("www.javatpoint.com", "javatpoint.com"),
		colly.Async(true),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		nextLink := e.Attr("href")
		e.Request.Visit(nextLink)
	})

	c.OnRequest(func(r *colly.Request) {
		// mcq url found
		// TODO: replace with regular expression
		if strings.Contains(r.URL.String(), "mcq") {
			fmt.Println(r.URL)
		}
	})
	c.Visit(baseUrl)
	c.Wait()
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
