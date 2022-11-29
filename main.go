package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"
	"RacoonApp/middleware"

	"github.com/MartijnKorbee/GoRacoon"
)

type application struct {
	App        *GoRacoon.GoRacoon
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
}

func main() {
	racoon := initApplication()
	racoon.App.ListenAndServe()
}
