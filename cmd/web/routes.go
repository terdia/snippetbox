package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	mux.Use(app.recoverPanic, app.logRequest, secureHeaders)

	mux.With(app.session.Enable).Get("/", app.home)
	mux.With(app.session.Enable).Get("/snippet/create", app.createSnippetForm)
	mux.With(app.session.Enable).Post("/snippet/create", app.createSnippet)
	mux.With(app.session.Enable).Get("/snippet/{id}", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
