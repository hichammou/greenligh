package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.hichammou/internal/data"
	"greenlight.hichammou/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a new anonymos struct to hold the information that we expect to be in the HTTP request body
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	//Initialize a new json.Decoder instance which reads from the request body and then use Decode() to decode the body contents into input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Init a new Validator instance
	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "year must be provided")
	v.Check(input.Year >= 1888, "year", "year must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "year must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Note that we're using the Unique helper in the line below to check that all
	// values in the input.Genres slice are unique.
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	fmt.Println(v.Errors)

	// Check if the fields passed the check
	if !v.Valide() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
