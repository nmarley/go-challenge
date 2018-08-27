package src_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func FormHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		res.WriteHeader(500)
		return
	}

	name := req.Form.Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(res, "Posted Hello, %s!", name)
}

func Test_U_FormHandler(t *testing.T) {
	// The url.Values type allows us to assemble a "form" that we can send as part
	// of the request.
	form := url.Values{}
	form.Add("name", "Ringo")

	// The `Encode` method on `url.Values` will properly encode the values we set
	// into well formed `string` that can be read as the body of the request.
	req := httptest.NewRequest("POST", "/form", strings.NewReader(form.Encode()))

	// We must set the `Content-Type` correctly for `ParseForm` to work.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	FormHandler(res, req)

	if got, exp := res.Code, http.StatusOK; got != exp {
		t.Errorf("unexpected response code.  got: %d, exp %d\n", got, exp)
	}
	if got, exp := res.Body.String(), "Posted Hello, Ringo!"; got != exp {
		t.Errorf("unexpected body.  got: %s, exp %s\n", got, exp)
	}
}

func Test_U_FormHandler_Template_Error(t *testing.T) {
	// pass invalid hex strings
	req := httptest.NewRequest("POST", "/form", strings.NewReader("%zzzzz"))

	// set the header `Content-Type` to `application/x-www-form-urlencoded`
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create a new httptest.NewRecorder
	res := httptest.NewRecorder()

	// Call the FormHandler
	FormHandler(res, req)

	// Test to see the the response code is 500
	if got, exp := res.Code, http.StatusInternalServerError; got != exp {
		t.Errorf("unexpected response code.  got: %d, exp %d\n", got, exp)
	}

	// test the body is `Oops!`
	if got, exp := res.Body.String(), "Oops!"; got != exp {
		t.Errorf("unexpected body.  got: %s, exp %s\n", got, exp)
	}
}
