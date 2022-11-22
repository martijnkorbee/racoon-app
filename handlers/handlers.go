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
	h.App.SessionManager.Put(r.Context(), "bar", "foo")

	err := h.App.Render.Page(w, r, "sessions", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}
