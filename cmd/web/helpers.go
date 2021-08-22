package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.logger.Error.Output(2, trace)

	_ = sentry.CaptureException(err)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	td.IsAuthenticated = app.isAuthenticated(r)

	// Add the flash message to the template data, if one exists.
	td.Flash = app.session.PopString(r, "flash")

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl').
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter.
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}

//todo: proper sentry implementation
func (app *application) initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://837f2bf7b3854cdaa07c2315de0293ec@o959477.ingest.sentry.io/5907807",
		//prod:Dsn: "https://f026770da7a24f5d908f6169f1540617@o959477.ingest.sentry.io/5908771",
		Environment:      "dev",
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		app.logger.Error.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
}
