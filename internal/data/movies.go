package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // the - tells the encoder to not show this field in the generated JSON
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`           // omitempty is used to tell the encoder to not show this field in the final JSON - if the value of this field is empty
	Runtime   Runtime   `json:"runtime,omitempty,string"` // string is to show this field as a string in JSON
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
