package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/terdia/snippetbox/pkg/forms"
	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
	"github.com/terdia/snippetbox/pkg/services"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	service := services.NewSnippetService(repository.NewSnippetRepository(app.DB))

	snippets, err := service.GetLatest()
	if err != nil {
		app.serverError(w, err)

		return
	}

	data := &templateData{Snippets: snippets}

	app.render(w, r, "home.page.tmpl", data)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)

		return
	}

	service := services.NewSnippetService(repository.NewSnippetRepository(app.DB))

	s, err := service.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)

		return
	}

	service := services.NewSnippetService(repository.NewSnippetRepository(app.DB))

	id, form, err := service.CreateSnippet(forms.New(r.PostForm))

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if form != nil {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})

		return
	}

	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
