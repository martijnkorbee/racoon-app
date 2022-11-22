package handlers

import (
	"net/http"

	"github.com/MartijnKorbee/GoRacoon"
)

type Handlers struct {
	App *GoRacoon.GoRacoon
}

// Home is the home route
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// Sessions is the sessions route
func (h *Handlers) SessionsTest(w http.ResponseWriter, r *http.Request) {
	myData := "foo"

	h.App.Session.Put(r.Context(), "bar", myData)

	err := h.App.Render.Page(w, r, "sessions", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}
