package render

import (
	"net/http"
	"testing"

	"github.com/Sri2103/bookings/internal/models"
)

func TestADdDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "12345")

	result := AddDefaultData(&td, r)

	if result.Flash != "12345" {
		t.Error("flash value of 12345 not found in the session")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, r, "home.page.html", &models.TemplateData{})

	if err != nil {
		t.Error("error writing template to browser")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("Get", "/some-url", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil

}

func TestNewTemplates(t *testing.T) {

	NewRenderer(app)

}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}
}
