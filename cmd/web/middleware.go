package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
	"github.com/terdia/snippetbox/pkg/services"
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if !app.isAuthenticated(r) {
			http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		rw.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(rw, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(
		http.Cookie{HttpOnly: true, Path: "/", Secure: true},
	)

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(rw, r)

			return
		}

		//to be cleanup up, during dependency injection refactoring
		userService := services.NewUserService(
			repository.NewUserRepository(app.DB),
			services.NewPasswordService(),
		)

		user, err := userService.GetById(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(rw, r)

			return
		} else if err != nil {
			app.serverError(rw, err)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)

		next.ServeHTTP(rw, r.WithContext(ctx))

	})
}
