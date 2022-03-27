package main

import (
	"goquiz/pkg/scrapper"
)

func main() {
	s := scrapper.New()
	s.GetQuizzes()
}
