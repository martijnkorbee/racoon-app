package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"
	"RacoonApp/middleware"

	"github.com/MartijnKorbee/GoRacoon"
)

type application struct {
	Racoon     *GoRacoon.GoRacoon
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
}

func main() {
	app := initApplication()
	app.Racoon.ListenAndServe()
}
