package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sri2103/bookings/internal/driver"
	"github.com/Sri2103/bookings/internal/helpers"

	"github.com/Sri2103/bookings/internal/config"
	"github.com/Sri2103/bookings/internal/handlers"
	"github.com/Sri2103/bookings/internal/models"
	"github.com/Sri2103/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLogger *log.Logger
var errorLogger *log.Logger

// main is the main function
func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Println("Starting mail listener.....")

	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.User{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)

	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false
	infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app.InfoLog = infoLogger

	errorLogger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLogger

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to a Database
	log.Println("connecting to database.....")
	db, err := driver.ConnectSQL("host = localhost port=5432 dbname=bookings user=postgres password=harsha")
	if err != nil {
		log.Fatal("Cannot connect to the Database! Dying....")
	}

	log.Println("connected to database.....")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
