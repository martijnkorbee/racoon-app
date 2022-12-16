package middleware

import (
	"RacoonApp/data"

	"github.com/martijnkorbee/goracoon"
)

type Middleware struct {
	App    *GoRacoon.GoRacoon
	Models *data.Models
}
