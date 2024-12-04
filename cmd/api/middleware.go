package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// If there was panic set a "Connection: Close" header on the response
				// This acts as a trigger to make Go's http server automatically close the current connection after a response has been sent
				w.Header().Set("Connection", "close")

				// the value of returned err is any. so we user fmt.Errorf() to normalize it into an error
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
