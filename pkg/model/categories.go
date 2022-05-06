package model

func (m *CategoriesModel) GetAll() ([]*Category, error) {

	return m.Categories, nil
}

func (m *CategoriesModel) GetByID(categId int) (*Category, error) {

	var category Category

	return &category, nil
}

func (m *CategoriesModel) Insert(categ Category) (int, error) {
	categ.Id = len(m.Categories)
	m.Categories = append(m.Categories, &categ)
	return categ.Id, nil
}
