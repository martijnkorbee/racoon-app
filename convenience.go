package main

import "net/http"

func (a *application) get(s string, h http.HandlerFunc) {
	a.Racoon.Routes.Get(s, h)
}

func (a *application) post(s string, h http.HandlerFunc) {
	a.Racoon.Routes.Post(s, h)
}

func (a *application) use(m ...func(http.Handler) http.Handler) {
	a.Racoon.Routes.Use(m...)
}
