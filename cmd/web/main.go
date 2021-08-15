package main

import (
	"flag"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/terdia/snippetbox/pkg/datasource"
	"github.com/terdia/snippetbox/pkg/logger"
	"github.com/terdia/snippetbox/pkg/models/mysql"
)

type application struct {
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	logger        *logger.SnippetLogger
}

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

	app := &application{
		snippets:      &mysql.SnippetModel{DB: connectionPool.DB},
		templateCache: templateCache,
		logger:        sLogger,
	}

	httpServer := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logger.Error,
		Handler:  app.routes(),
	}

	app.logger.Info.Printf("Starting server on %s", *addr)
	err = httpServer.ListenAndServe()

	app.logger.Error.Fatal(err)
}
