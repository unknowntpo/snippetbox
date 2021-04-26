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

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestAbout(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	wantCode := http.StatusOK
	aboutUrl := "/about"
	wantBody := []byte(`A great about page.`)

	code, _, body := ts.get(t, aboutUrl)
	if code != wantCode {
		t.Errorf("want %d; got %d", wantCode, code)
	}

	if !bytes.Contains(body, wantBody) {
		t.Errorf("got %q, want body to contain %q", body, wantBody)
	}
}

// TestCreateSnippetForm tests the CreateSnippetForm() which should have these behavior:
// 1. Unauthenticated users are redirected to the login form.
// 2. Authenticated users are redirected to the login form.
func TestCreateSnippetForm(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, body := ts.get(t, "/snippet/create")
		wantCode := http.StatusSeeOther

		if code != wantCode {
			t.Errorf("want %v, got %v", wantCode, code)
		}
		// Check body
		gotHeader := headers.Get("Location")
		wantHeader := "/user/login"
		if gotHeader != wantHeader {
			t.Errorf("want header to be %v, got %v", wantHeader, gotHeader)
		}

		const seeOther = `<a href="/user/login">See Other</a>`
		wantBody := []byte(seeOther)

		if !bytes.Contains(body, wantBody) {
			t.Errorf("got %q, want body to contain %q", body, wantBody)
		}
	})
	t.Run("Authenticated", func(t *testing.T) {
		// GET /snippet/login
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		// use email: "alice@example.com", passwd: ""
		form := url.Values{}
		form.Add("name", "alice")
		form.Add("email", "alice@example.com")
		form.Add("password", "")
		form.Add("csrf_token", csrfToken)
		code, _, body := ts.postForm(t, "/user/login", form)
		// check if login succeed
		wantCode := http.StatusSeeOther
		if code != wantCode {
			t.Errorf("Wrong code, want %v, got %v", wantCode, code)
		}

		// GET /snippet/create
		code, _, body = ts.get(t, "/snippet/create")
		wantCode = http.StatusOK
		if code != wantCode {
			t.Errorf("Wrong code, want %v, got %v", wantCode, code)
		}

		// Verify content : <form action='/snippet/create' method='POST'>
		formTag := `<form action='/snippet/create' method='POST'>`
		if !bytes.Contains(body, []byte(formTag)) {
			t.Errorf("got %q, want body to contain %q", body, formTag)
		}
	})
}
