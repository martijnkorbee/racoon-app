package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"
	"RacoonApp/middleware"
	"log"
	"os"

	"github.com/martijnkorbee/goracoon"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init GoRacoon
	racoon := &goracoon.Goracoon{}
	err = racoon.New(path)
	if err != nil {
		log.Fatal(err)
	}

	racoon.AppName = "RacoonApp"

	myMiddleware := &middleware.Middleware{
		App: racoon,
	}

	myHandlers := &handlers.Handlers{
		App: racoon,
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
