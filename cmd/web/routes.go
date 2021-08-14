package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	//standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := chi.NewRouter()

	mux.Use(app.recoverPanic, app.logRequest, secureHeaders)

	mux.Get("/", app.home)
	mux.Get("/snippet/create", app.createSnippetForm)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Get("/snippet/{id}", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
