package main

import (
	"RacoonApp/handlers"

	"github.com/MartijnKorbee/GoRacoon"
)

type application struct {
	App      *GoRacoon.GoRacoon
	Handlers *handlers.Handlers
}

func main() {
	racoon := initApplication()
	racoon.App.ListenAndServe()
}
