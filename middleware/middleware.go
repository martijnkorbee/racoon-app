package middleware

import (
	"RacoonApp/data"

	"github.com/MartijnKorbee/GoRacoon"
)

type Middleware struct {
	App    *GoRacoon.GoRacoon
	Models *data.Models
}
