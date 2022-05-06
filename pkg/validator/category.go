package validator

import (
	"fmt"
	"goquiz/pkg/model"
)

func (v *Validator) ValidateCategoryId(categId int) {
	v.Check(categId > 0, "category_id", "must greater than zero")
	// TODO: dynamic categ no based on the numbers of categories in the database
	v.Check(categId < 200, "category_id", "must be less than 200")

}

func (v *Validator) ValidateCategory(categ *model.Category) {
	v.Check(len(categ.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(Matches(categ.Name, NonSpecialText), "name", "name must not contain any special characters")
}

func (v *Validator) ValidateQuestion(que *model.Question) {
	v.Check(que.Text != "", "text", "must not be empty")
	v.Check(len(que.Text) <= 1000, "text", "must not be more than 1000 bytes long")

	v.Check(que.Category.Id > 0, "category_id", "must be greater than zero")
	// TODO: dynamic categ no based on the max number of categ ID in the database
	v.Check(que.Category.Id < 500, "category_id", "must be less than 500")

	v.Check(len(que.AnsOptions) > 2, "answers", "must be more than 2 options")
	v.Check(len(que.AnsOptions) <= 5, "answers", "must not be more than 5 options")
	for _, ans := range que.AnsOptions {
		v.Check(ans != "", "answers", "must not contain empty answer")
	}

	// TODO: refactor que.Answer.Option+1.
	v.Check(len(que.Answer.Explanation) <= 2000, "explanation", "must not be more than 2000 bytes long")
	v.Check(que.Answer.Option+1 > 0, "correct_answer", "must be more than 0")
	v.Check(int(que.Answer.Option) < len(que.AnsOptions), "correct_answer",
		fmt.Sprintf("must be within answers range. must not be more than %d", len(que.AnsOptions)),
	)

	v.Check(len(que.Codeblock) <= 2000, "codeblock", "must not be more than 2000 bytes long")
}
