package middleware

import (
	"racoonapp/data"
)

type Middleware struct {
	App    *goracoon.goracoon
	Models *data.Models
}
