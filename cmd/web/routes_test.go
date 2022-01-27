package main

import (
	"fmt"
	"testing"

	"github.com/Sri2103/bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRouttes(t *testing.T) {

	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {

	case *chi.Mux:
		//do nothing

	default:
		t.Error(fmt.Printf("this is not *chi.Mux, but it is %T", v))

	}

}
