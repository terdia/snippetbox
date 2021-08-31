package main

import (
	"net/http"

	"github.com/terdia/snippetbox/pkg/forms"
)

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form, err := app.userService.SignupUser(forms.New(r.PostForm))

	if form != nil {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})

		return
	}

	if err != nil {
		app.serverError(w, err)

		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	authResponse := app.userService.Authenticate(forms.New(r.PostForm))
	form := authResponse.Form
	err = authResponse.Error
	user := authResponse.User

	if form != nil && !form.Valid() {
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "authenticatedUserID", user.ID)

	redirectPathAfterLogin := app.session.PopString(r, "redirectPathAfterLogin")
	if redirectPathAfterLogin == "" {
		redirectPathAfterLogin = "/snippet/create"
	}

	http.Redirect(w, r, redirectPathAfterLogin, http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userProfle(w http.ResponseWriter, r *http.Request) {

	user, err := app.userService.GetById(app.session.GetInt(r, "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.render(w, r, "profile.page.tmpl", &templateData{User: user})
}
