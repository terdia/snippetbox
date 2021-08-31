package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	rs := ts.get(t, "/ping")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	if string(rs.Body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := ts.get(t, test.urlPath)

			if res.StatusCode != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, res.StatusCode)
			}

			if !bytes.Contains(res.Body, test.wantBody) {
				t.Errorf("want body to contain %q", test.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, "/user/signup")

	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, res.StatusCode)
	}

	csrfToken := extractCSRFToken(t, res.Body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Terry", "terry@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "terry@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is required.")},
		{"Empty email", "Terry", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is required.")},
		{"Empty password", "Terry", "terry@example.com", "", csrfToken, http.StatusOK, []byte("This field is required.")},
		{"Invalid email (incomplete domain)", "Terry", "Terry@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Terry", "Terryexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Terry", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Terry", "terry@example.com", "pa$$w", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 6 characters)")},
		{"Duplicate email", "Terry", "duplicate@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			res := ts.postForm(t, "/user/signup", form)

			if res.StatusCode != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, res.StatusCode)
			}

			if !bytes.Contains(res.Body, test.wantBody) {
				t.Errorf("want body %s to contain %q", res.Body, test.wantBody)
			}
		})
	}
}

//Unauthenticated users are redirected to the login form.
//Authenticated users are shown the form to create a new snippet.

func TestCreateSnippetForm(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated User", func(t *testing.T) {
		res := ts.get(t, "/snippet/create")

		if res.StatusCode != http.StatusSeeOther {
			t.Errorf("wants %d; got %d", http.StatusSeeOther, res.StatusCode)
		}

		location := res.Header.Get("Location")

		if location != "/user/login" {
			t.Errorf("wants %s; got %s", "/user/login", location)
		}
	})

	t.Run("Authenticated User", func(t *testing.T) {
		loginPageRes := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, loginPageRes.Body)

		loginForm := url.Values{}
		loginForm.Add("email", "terry@yahoo.com")
		loginForm.Add("password", "1234567")
		loginForm.Add("csrf_token", csrfToken)

		_ = ts.postForm(t, "/user/login", loginForm)

		resp := ts.get(t, "/snippet/create")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("wants %d; got %d", http.StatusOK, resp.StatusCode)
		}

		formTag := "<form action='/snippet/create' method='POST'>"
		if !bytes.Contains(resp.Body, []byte(formTag)) {
			t.Errorf("want body %s to contain %q", resp.Body, formTag)
		}

	})

}
