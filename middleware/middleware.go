package middleware

import (
	"racoonapp/data"

	"github.com/martijnkorbee/goracoon"
)

type Middleware struct {
	Racoon *goracoon.Goracoon
	Models *data.Models
}
