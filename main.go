package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"
	"RacoonApp/middleware"

	"github.com/martijnkorbee/goracoon"
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
