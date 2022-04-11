package validator

func (v *Validator) ValidateCategoryId(categId int) {
	v.Check(categId > 0, "category_id", "must greater than zero")
	v.Check(categId < 1000, "category_id", "must be less than 100")
}
