package postgres

import "database/sql"

type Model struct {
	QuestionsModel
	CategoriesModel
}

type QuestionsModel struct {
	DB *sql.DB
}

type CategoriesModel struct {
	DB *sql.DB
}
