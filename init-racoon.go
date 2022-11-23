package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"
	"log"
	"os"

	"github.com/MartijnKorbee/GoRacoon"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init GoRacoon
	racoon := &GoRacoon.GoRacoon{}
	err = racoon.New(path)
	if err != nil {
		log.Fatal(err)
	}

	racoon.AppName = "RacoonApp"

	myHandlers := &handlers.Handlers{
		App: racoon,
	}

	app := &application{
		App:      racoon,
		Handlers: myHandlers,
	}

	app.App.Routes = app.routes()

	app.Models = data.New(app.App.DB.ConnectionPool)

	return app
}
