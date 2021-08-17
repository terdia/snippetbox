package session

import (
	"os"
	"time"

	"github.com/golangcollege/sessions"
)

const (
	appSessionSecret = "SESSION_SECRET"
)

var (
	sessionSecret = os.Getenv(appSessionSecret)
)

func NewSession() *sessions.Session {
	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(sessionSecret))
	session.Lifetime = 12 * time.Hour

	return session
}
