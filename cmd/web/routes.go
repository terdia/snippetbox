package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type dynamicMiddlewares []func(http.Handler) http.Handler

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	router.Use(app.recoverPanic, app.logRequest, secureHeaders)

	dm := dynamicMiddlewares{
		app.session.Enable,
		noSurf,
		app.authenticate,
	}

	router.With(dm...).Get("/", app.home)
	router.With(append(dm, app.requireAuthentication)...).Get("/snippet/create", app.createSnippetForm)
	router.With(append(dm, app.requireAuthentication)...).Post("/snippet/create", app.createSnippet)
	router.With(dm...).Get("/snippet/{id}", app.showSnippet)

	router.With(dm...).Get("/user/signup", app.signupUserForm)
	router.With(dm...).Post("/user/signup", app.signupUser)
	router.With(dm...).Get("/user/login", app.loginUserForm)
	router.With(dm...).Post("/user/login", app.loginUser)
	router.With(append(dm, app.requireAuthentication)...).Post("/user/logout", app.logoutUser)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
