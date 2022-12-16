package handlers

import (
	"RacoonApp/data"
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/martijnkorbee/goracoon"
)

type Handlers struct {
	App    *goracoon.Goracoon
	Models *data.Models
}

// Home is the home route
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// Sessions is the sessions route
func (h *Handlers) SessionsTest(w http.ResponseWriter, r *http.Request) {
	// create session data
	h.sessionPut(r.Context(), "foo", "bar")
	// extract session data
	sessionData := h.sessionGet(r.Context(), "foo")

	// create jet varmap
	vars := make(jet.VarMap)
	// add session data
	vars.Set("foo", sessionData)

	err := h.render(w, r, "sessions", vars, nil)
	if err != nil {
		h.logError("error rendering", err)
	}
}

func (h *Handlers) TestCrypto(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello, World!"
	fmt.Fprintln(w, "Unencrypted: "+plainText)

	encrypted, err := h.encrypt(plainText)
	if err != nil {
		h.logError("error encrypting: ", err)
		h.App.Error500(w, r)
		return
	}

	fmt.Fprintln(w, "Encrypted: "+encrypted)

	decrypted, err := h.decrypt(encrypted)
	if err != nil {
		h.logError("error decrypting: ", err)
		h.App.Error500(w, r)
		return
	}

	fmt.Fprintln(w, "Decrypted: "+decrypted)
}
