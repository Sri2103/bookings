package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	r := httptest.NewRequest("Post", "/whatever", nil)

	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when it should have been valid")
	}

}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("Post", "/whatever", nil)

	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("got valid when it should have missing values")
	}

	postedData := url.Values{}

	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("Post", "/whatever", nil)

	r.PostForm = postedData

	form = New(r.PostForm)

	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("shows does not have required values when it does")
	}

}

func TestForm_Has(t *testing.T) {

	r := httptest.NewRequest("Post", "/whatever", nil)

	form := New(r.PostForm)

	has := form.Has("/whatever")

	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}

	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")

	if !has {
		t.Error("returned empty when there is s omething already in form")
	}

}

func TestForm_MinLength(t *testing.T) {

	r := httptest.NewRequest("Post", "/whatever", nil)

	form := New(r.PostForm)

	form.MinLength("x", 10)

	if form.Valid() {
		t.Error("form minLength for non existent fields")
	}

	isError := form.Errors.Get("x")

	if isError == "" {
		t.Error("should have an error but did not get one.")
	}

	postedValue := url.Values{}

	postedValue.Add("some_field", "some field")
	form = New(postedValue)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows the length of 100 is met where it is shorter")
	}

	postedValue = url.Values{}

	postedValue.Add("other_field", "abc123")
	form = New(postedValue)

	form.MinLength("other_field", 1)

	isError = form.Errors.Get("other_field")

	if isError != "" {
		t.Error("should not have an error but did  get one.")
	}

	if !form.Valid() {
		t.Error("shows the length of 1 is not met where it is ")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}

	form := New(postedValues)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("form show valid email for non-existent field")
	}

	postedValues = url.Values{}

	postedValues.Add("email", "me@here.com")

	form = New(postedValues)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("got an invalid Email when we should not have")
	}

	postedValues = url.Values{}

	postedValues.Add("email", "x")

	form = New(postedValues)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("got a valid Email when we should not have")
	}

}
