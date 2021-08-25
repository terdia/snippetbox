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

	router.With(app.session.Enable, noSurf).Get("/", app.home)
	router.With(app.session.Enable, app.requireAuthentication, noSurf).Get("/snippet/create", app.createSnippetForm)
	router.With(app.session.Enable, app.requireAuthentication, noSurf).Post("/snippet/create", app.createSnippet)
	router.With(app.session.Enable, noSurf).Get("/snippet/{id}", app.showSnippet)

	router.With(app.session.Enable, noSurf).Get("/user/signup", app.signupUserForm)
	router.With(app.session.Enable, noSurf).Post("/user/signup", app.signupUser)
	router.With(app.session.Enable, noSurf).Get("/user/login", app.loginUserForm)
	router.With(app.session.Enable, noSurf).Post("/user/login", app.loginUser)
	router.With(app.session.Enable, app.requireAuthentication, noSurf).Post("/user/logout", app.logoutUser)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
