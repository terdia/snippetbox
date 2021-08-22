package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	router.Use(app.recoverPanic, app.logRequest, secureHeaders)

	router.With(app.session.Enable).Get("/", app.home)
	router.With(app.session.Enable, app.requireAuthentication).Get("/snippet/create", app.createSnippetForm)
	router.With(app.session.Enable, app.requireAuthentication).Post("/snippet/create", app.createSnippet)
	router.With(app.session.Enable).Get("/snippet/{id}", app.showSnippet)

	router.With(app.session.Enable).Get("/user/signup", app.signupUserForm)
	router.With(app.session.Enable).Post("/user/signup", app.signupUser)
	router.With(app.session.Enable).Get("/user/login", app.loginUserForm)
	router.With(app.session.Enable).Post("/user/login", app.loginUser)
	router.With(app.session.Enable, app.requireAuthentication).Post("/user/logout", app.logoutUser)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
