package main

import (
	"racoonapp/data"
	"racoonapp/handlers"
	"racoonapp/middleware"

	"github.com/martijnkorbee/goracoon"
)

type application struct {
	Racoon     *goracoon.Goracoon
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
}

func main() {
	app := initApplication()
	app.Racoon.ListenAndServe()
}
