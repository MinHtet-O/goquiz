package postgres

import "database/sql"

type Model struct {
	Questions  Questions
	Categories Categories
}

type Questions struct {
	DB *sql.DB
}

type Categories struct {
	DB *sql.DB
}
