package model

func (m CategoriesModel) GetAll() ([]*Category, error) {

	return m.Categories, nil
}

func (m CategoriesModel) GetByID(categId int) (*Category, error) {

	var category Category

	return &category, nil
}

func (m CategoriesModel) Insert(cate Category) (int, error) {
	return 0, nil
}
