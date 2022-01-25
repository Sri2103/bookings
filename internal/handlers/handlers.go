package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Sri2103/bookings/internal/config"
	"github.com/Sri2103/bookings/internal/forms"
	"github.com/Sri2103/bookings/internal/models"
	"github.com/Sri2103/bookings/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})

}

// renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})

}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{})

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJson handles request for availability and sends Json response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// perform some logic

	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)

	form.Required("first_name", "last_name", "email", "phone")

	form.MinLen("first_name", 3, r)

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Data: data,
			Form: form,
		})

		return
	}

}
