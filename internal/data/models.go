package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound err.
var (
	ErrRecordNotFound = errors.New("")
)

type Models struct {
	Movies MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
