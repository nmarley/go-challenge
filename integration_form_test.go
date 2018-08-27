package src_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func App() http.Handler {
	mux := http.NewServeMux()
	// mux.HandleFunc("/hello", HelloHandler) // GET
	mux.HandleFunc("/form", FormHandler) // POST

	h := func(res http.ResponseWriter, req *http.Request) {
		rs := &Responder{ResponseWriter: res, status: 200}
		mux.ServeHTTP(rs, req)
		if rs.status == 500 {
			rs.Write([]byte("Oops!"))
		}
	}
	return http.HandlerFunc(h)
}

func Test_I_FormHandler_Error_Template(t *testing.T) {
	// Create a new httptest.NewServer with our App
	ts = httptest.NewServer(App())
	// Defer the test server Close
	defer ts.Close()

	// Post `%zzzz` to the `/form` endpoint
	req, err := http.Post(ts.URL+"/form",
		"application/x-www-form-urlencoded",
		strings.NewReader("%zzzz"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Test the status code is http.StatusInternalServerError
	if got, exp := res.StatusCode, http.StatusInternalServerError; got != exp {
		t.Error("unexpected status code: got %d, expected %d", got, exp)
	}

	b, err := ioutil.ReadAll(res.body)
	if err != nil {
		t.Fatal(err)
	}

	// test the body is `Oops!`
	if got, exp := string(b), "Oops!"; got != exp {
		t.Error("unexpected body: got %s, expected %s", got, exp)
	}
}
