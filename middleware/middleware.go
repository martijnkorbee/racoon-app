package middleware

import (
	"racoonapp/data"

	"github.com/martijnkorbee/goracoon"
)

type Middleware struct {
	App    *goracoon.Goracoon
	Models *data.Models
}
