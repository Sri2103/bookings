package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

//New initializes a form struct
func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}

}

// Returns true if there are no errors and false if ther are errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "the field cannot be blank")
		}
	}
}

//Has func will check if form field is in Post and is not empty
func (f *Form) Has(field string, r *http.Request) bool {

	x := r.Form.Get(field)

	if x == "" {
		f.Errors.Add(field, "this field cannot be empty ")
		return false
	}
	return true

}

func (f *Form) MinLen(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("this field must be atleast %d characters long", length))
		return false
	}
	return true
}
