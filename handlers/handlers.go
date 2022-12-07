package handlers

import (
	"net/http"
	"racoonapp/data"
)

type Handlers struct {
	App    *goracoon.goracoon
	Models *data.Models
}

// Home is the home route
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
