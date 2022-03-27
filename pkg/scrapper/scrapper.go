package scrapper

// TODO: refactor for proper naming and package
import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	model "goquiz/pkg/model"
	"goquiz/utils"
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

func New() *QuizScrapper {
	rootDomain := "javatpoint.com"
	filepath := "./mcq_url.txt"
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
			utils.SaveQuestionToFile(s.rootDomain, categ, fmt.Sprintf("%v", questions))
			if len(*questions) > minlimit {
				s.mu.Lock()
				s.quizzes.AddQuestions(categ, *questions)
				s.mu.Unlock()
			} else {
				fmt.Fprintln(os.Stderr, "Remove ", categ, " from questions with length", len(*questions))
			}
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
		ansReg  = regexp.MustCompile(`(\(\w\))|(\s\w\.$)|(\s\w\)\s)|(\s\w$)|(\s\w\.\s\w*)|(\s\w:\s)|(\s\w\s)`)
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

		// post with second 'pq' text which need to append into question text
		if quesNode := utils.FindSibling(e.DOM, "pq", 4, utils.D_Next); quesNode != nil {
			question.Text += "\n" + quesNode.Text()
		}

		// ignore second 'pq' text for the post
		if quesNode := utils.FindSibling(e.DOM, "pq", 4, utils.D_Prev); quesNode != nil {
			fmt.Fprintf(os.Stderr, "Second pq question text found category %s , index %d \n", category, e.Index)
			return
		}

		// post that contains codeblock, codeblock is optional
		if codeblockNode := utils.FindSibling(e.DOM, "codeblock", 4, utils.D_Next); codeblockNode != nil {
			question.Codeblock = codeblockNode.Text()
		}

		// answer options as text
		// valid A,B,C,D,E options
		func() {
			if optsNode := utils.FindSibling(e.DOM, "pointsa", 6, utils.D_Next); optsNode != nil && optsNode.Children().Length() < 6 {
				optsNode.Children().Each(func(i int, c *goquery.Selection) {
					question.Options[model.Option(i)] = strings.ToLower(c.Text())
				})
				// anser options as image
			} else if imageURL, exist := e.DOM.Next().Children().First().Attr("src"); exist {
				err := utils.ValidateImgURL(imageURL)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid Image URL", imageURL, " for category ", category)
					return
				}
				question.Options[0] = imageURL

			} else { // invalid answer option
				fmt.Fprintln(os.Stderr, "Invalid Answer Options ", category, ", index", e.Index)
				return
			}
		}()
		// correct ans for the question
		if ansNode := utils.FindSibling(e.DOM, "testanswer", 8, utils.D_Next); ansNode != nil {
			// index 0 is the correct answer, 1 to ... is explanation
			ansNode.Children().Each(func(i int, c *goquery.Selection) {
				if i == 0 {
					// extract the ans option from the whole string
					ans := ansReg.FindString(c.Text())
					if ans != "" {
						ans = strings.Split(ans, "")[1]
						ans = strings.ToLower(ans)
					}
					opt, found := model.AnsMapping[ans]
					// set the correct answer
					if found {
						question.CorrectAns.Option = opt
					}
				} else {
					// append the explanation
					question.CorrectAns.Explanation += c.Text()
				}
			})
		} else {

			fmt.Fprintln(os.Stderr, "Invalid Correct Answer ", category, ", index", e.Index)
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
		} else {
			if r == ')' {
				breakIndex = i
				break
			}
			valid = false
			break
		}
	}
	if valid {
		title = title[breakIndex+2:]
	}
	return title
}
