package main

import (
	"RacoonApp/data"
	"RacoonApp/handlers"

	"github.com/MartijnKorbee/GoRacoon"
)

type application struct {
	App      *GoRacoon.GoRacoon
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	racoon := initApplication()
	racoon.App.ListenAndServe()
}
