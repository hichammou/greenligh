package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// If there was panic set a "Connection: Close" header on the response
				// This acts as a trigger to make Go's http server automatically close the current connection after a response has been sent
				w.Header().Set("Connection", "close")

				// The value of returned err is any. so we user fmt.Errorf() to normalize it into an error
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the client's IP addresses and rate limiters
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Lunch a background cleanup which removes old entries one every minute
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any data race
			mu.Lock()

			// Remove any ip that was not sent any request in the last 3 minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// Unlock the mutex
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limiter is enabled.
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()

			if _, found := clients[ip]; !found {
				// Create and add a new client struct to the map if it doesn't already exists.
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// Update the last seen time for the client
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExccededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
