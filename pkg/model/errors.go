package model

import "errors"

// datbase related errors
var (
	ErrRecordNotFound = errors.New("record not found")
)

// validaton errors

var (
	MustProvideInt = "must provide integer value"
)
