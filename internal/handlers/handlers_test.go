package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"Home", "/", "GET", []postData{}, http.StatusOK},
	{"About", "/about", "GET", []postData{}, http.StatusOK},
	{"Gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"Ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"Sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"MakeRes", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},

	{"possearch-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Harsha"},
		{key: "last_name", value: "Dama"},
		{key: "email", value: "h@h.com"},
		{key: "phone", value: "412-555-1234"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTest {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d received %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {

			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d received %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}

	}
}
