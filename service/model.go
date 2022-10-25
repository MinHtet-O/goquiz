package service

type Model struct {
	QuestionsModel  QuestionsModel
	CategoriesModel CategoriesModel
}

type QuestionsModel interface {
	GetAll(category Category) ([]Question, error)
	Insert(categID int, q Question) (int, error)
}

type CategoriesModel interface {
	GetAll() ([]*Category, error)
	GetByID(categId int) (*Category, error)
	Insert(cate Category) (int, error)
}

type Questions []Question

type Option int8

const (
	A Option = iota
	B
	C
	D
	E
	O_MAX
)

var AnsMapping = map[string]Option{
	"a": A,
	"b": B,
	"c": C,
	"d": D,
	"e": E,
}

type Category struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	Questions      []Question `json:"-"`
	QuestionsCount int32      `json:"questions_count,omitempty"`
}

type Question struct {
	Id       int    `json:"id"`
	WebIndex int    `json:"-"`
	Text     string `json:"text"`
	// TODO: change the name from AnswerOptions to Choices
	AnsOptions []string `json:"answers"`
	Codeblock  string   `json:",omitempty"`
	Answer     Answer   `json:"correct_answer"`
	URL        string   `json:"-"`
	// TODO: remove Category from Question
	Category Category `json:"-"`
}

type Answer struct {
	Option      Option `json:"option"`
	Explanation string `json:"explanation"`
}

/*
// These are data structures to implement quiz games, consul/ web game that allow user to choose categories
// and guess the answer, development is not in progress currently

type User struct {
	name string
}

// add method to change the setting - setters
type Setting struct {
	quesTimeout int
}

//MCQ run time data structure
// add method to calculate totalScore
type MatchRecord struct {
	choices    []Choice
	date       time.Time
	totalScore int
}

type Choice struct {
	question Question
	ans      Option
	duration time.Time
}
*/
