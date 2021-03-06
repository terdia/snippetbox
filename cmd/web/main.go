package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/terdia/snippetbox/pkg/datasource"
	"github.com/terdia/snippetbox/pkg/logger"
	"github.com/terdia/snippetbox/pkg/repository"
	"github.com/terdia/snippetbox/pkg/services"
	"github.com/terdia/snippetbox/pkg/session"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	logger         *logger.SnippetLogger
	session        *sessions.Session
	templateCache  map[string]*template.Template
	userService    services.UserServiceInterface
	snippetService services.SnippetServiceInterface
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
		templateCache: templateCache,
		logger:        sLogger,
		session:       session.NewSession(),
		userService: services.NewUserService(
			repository.NewUserRepository(connectionPool.DB),
			services.NewPasswordService(),
		),
		snippetService: services.NewSnippetService(repository.NewSnippetRepository(connectionPool.DB)),
	}

	app.initSentry()

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	httpServer := &http.Server{
		Addr:         *addr,
		ErrorLog:     app.logger.Error,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info.Printf("Starting server on %s", *addr)

	//todo fix getting cert
	cert := "./tls/cert.pem"
	certKey := "./tls/key.pem"
	err = httpServer.ListenAndServeTLS(cert, certKey)

	app.logger.Error.Fatal(err)
}
