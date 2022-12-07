package main

import (
	"racoonapp/data"
	"racoonapp/handlers"
	"racoonapp/middleware"

	_ "github.com/martijnkorbee/goracoon"
)

type application struct {
	Racoon     *goracoon.goracoon
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
}

func main() {
	app := initApplication()
	app.Racoon.ListenAndServe()
}
