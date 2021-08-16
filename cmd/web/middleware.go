package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// Any code here will execute on the way down the chain.
		rw.Header().Set("X-XSS-Protection", "1; mode=block")
		rw.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(rw, r)
		// Any code here will execute on the way back up the chain.
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//log every request like so e.g. 172.18.0.1:60504 - HTTP/1.1 GET /snippet?id=4
		app.logger.Info.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(rw, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Check if there has been a panic or not use the builtin recover function.
			if err := recover(); err != nil {
				//sentry.CurrentHub().Recover(err)
				//sentry.Flush(time.Second * 5)
				// Set a "Connection: close" header on the response,
				//a triggers an automatic close of current connection.
				//after a response has been sent
				rw.Header().Set("Connection", "close")

				app.serverError(rw, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
