package handlers

import (
	"net/http"
	"racoonapp/data"

	"github.com/martijnkorbee/goracoon"
)

type Handlers struct {
	Racoon *goracoon.Goracoon
	Models *data.Models
}

// Home is the home route
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Racoon.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.Log.Error().Err(err).Msg("error rendering")
	}
}
