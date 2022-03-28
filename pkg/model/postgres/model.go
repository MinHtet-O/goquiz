package postgres

import "database/sql"

type DBmodel struct {
	Quizzes QuizzModel
}

type QuizzModel struct {
	DB *sql.DB
}
