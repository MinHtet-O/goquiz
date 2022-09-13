package postgres

import (
	"database/sql"
)

type PostgresModel struct {
	QuestionsModel
	CategoriesModel
}

type QuestionsModel struct {
	DB *sql.DB
}

type CategoriesModel struct {
	DB *sql.DB
}
