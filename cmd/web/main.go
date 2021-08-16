package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/terdia/snippetbox/pkg/datasource"
	"github.com/terdia/snippetbox/pkg/logger"
	"github.com/terdia/snippetbox/pkg/models/mysql"
)

type application struct {
	logger        *logger.SnippetLogger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

const (
	appSessionSecret = "SESSION_SECRET"
)

var (
	sessionSecret = os.Getenv(appSessionSecret)
)

func main() {

	addr := flag.String("addr", ":4000", "Http network address.")
	flag.Parse()

	sLogger := logger.New()

	connectionPool, err := datasource.NewConnectionPool()
	if err != nil {
		sLogger.Error.Fatal(err)
	}
	defer connectionPool.DB.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		sLogger.Error.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(sessionSecret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		snippets:      &mysql.SnippetModel{DB: connectionPool.DB},
		templateCache: templateCache,
		logger:        sLogger,
		session:       session,
	}

	app.initSentry()

	httpServer := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logger.Error,
		Handler:  app.routes(),
	}

	app.logger.Info.Printf("Starting server on %s", *addr)
	err = httpServer.ListenAndServe()

	app.logger.Error.Fatal(err)
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
