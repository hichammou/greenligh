package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a new line at the end of json for terminal apps
	js = append(js, '\n')

	// maps.Insert(w.Header(), maps.All(headers))

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	// DisallowUnknownFields() tels the decoder to return an error if JSON in the request body contains any fields that can not be mapped to target destination
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)

	if err != nil {
		var (
			syntaxError            *json.SyntaxError
			unmarshalTypeError     *json.UnmarshalTypeError
			invalidUnmarshaldError *json.InvalidUnmarshalError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// We check if the decoder found a field that can not be mapped to dst, then extract the field name from the error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Check if the request exceeds 1MB
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger that %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshaldError):
			panic(err)

		default:
			return err
		}
	}

	// Here we called Decode() again to check if the request body contains only a single JSON value. if not we return an error
	err = dec.Decode(&struct{})
	if err != io.EOF {
		return errors.New("body must contain only one single JSON value")
	}

	return nil
}
