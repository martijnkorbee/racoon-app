package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
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
	// create session data
	h.App.SessionManager.Put(r.Context(), "foo", "bar")
	// extract session data
	sessionData := h.App.SessionManager.GetString(r.Context(), "foo")

	// create jet varmap
	vars := make(jet.VarMap)
	// add session data
	vars.Set("foo", sessionData)

	err := h.App.Render.Page(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}
