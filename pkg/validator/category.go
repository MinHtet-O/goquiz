package validator

import "goquiz/pkg/model"

func (v *Validator) ValidateCategoryId(categId int) {
	v.Check(categId > 0, "category_id", "must greater than zero")
	v.Check(categId < 1000, "category_id", "must be less than 100")

}

func (v *Validator) ValidateCategory(categ *model.Category) {
	v.Check(len(categ.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(Matches(categ.Name, NonSpecialText), "name", "name must not contain any special characters")
}
