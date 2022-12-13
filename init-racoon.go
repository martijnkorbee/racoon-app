package main

import (
	"log"
	"os"
	"racoonapp/data"
	"racoonapp/handlers"
	"racoonapp/middleware"

	"github.com/martijnkorbee/goracoon"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init goracoon
	racoon := &goracoon.Goracoon{}
	err = racoon.New(path)
	if err != nil {
		log.Fatal(err)
	}

	racoon.AppName = "racoonapp"

	myMiddleware := &middleware.Middleware{
		Racoon: racoon,
	}

	myHandlers := &handlers.Handlers{
		Racoon: racoon,
	}

	app := &application{
		Racoon:     racoon,
		Middleware: myMiddleware,
		Handlers:   myHandlers,
	}

	app.Racoon.Routes = app.routes()

	app.Models = data.New(app.Racoon.DB.ConnectionPool)
	myHandlers.Models = &app.Models
	app.Middleware.Models = &app.Models

	return app
}
