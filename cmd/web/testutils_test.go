package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/terdia/snippetbox/pkg/logger"
	"github.com/terdia/snippetbox/pkg/repository/mock"
	"github.com/terdia/snippetbox/pkg/services"
)

type IntegrationTestResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		templateCache: templateCache,
		session:       session,
		logger: &logger.SnippetLogger{
			Info:  log.New(io.Discard, "", 0),
			Error: log.New(io.Discard, "", 0),
		},
		userService:    services.NewUserService(mock.NewUserRepository(), services.NewPasswordService()),
		snippetService: services.NewSnippetService(mock.NewSnippetRepository()),
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the client, so that response cookies are stored
	// and then sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) Get(t *testing.T, urlPath string) IntegrationTestResponse {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return IntegrationTestResponse{
		StatusCode: rs.StatusCode,
		Header:     rs.Header,
		Body:       body,
	}
}
