package middleware

import (
	"RacoonApp/data"

	"github.com/martijnkorbee/goracoon"
)

type Middleware struct {
	App    *goracoon.Goracoon
	Models *data.Models
}
