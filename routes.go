package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.Racoon.Routes.Get("/", a.Handlers.Home)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.Racoon.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.Racoon.Routes
}
