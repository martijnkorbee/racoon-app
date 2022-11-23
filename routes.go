package main

import (
	"RacoonApp/data"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/sessions", a.Handlers.SessionsTest)

	// test routes
	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Martijn",
			LastName:  "Korbee",
			Email:     "test@test.nl",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.AddUser(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "User created with ID: %d", id)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
