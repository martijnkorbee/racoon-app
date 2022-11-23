package handlers

import "net/http"

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	// parse the form post
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	// extract form data
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// get user from db
	user, err := h.Models.Users.GetUserByEmail(email)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	// check if password matches
	matches, err := user.AuthenticateUser(password)
	if err != nil {
		h.App.ErrorLog.Println("Error validating password", err)
		return
	}
	if !matches {
		w.Write([]byte("Invalid password!"))
		return
	}

	// add a session
	h.App.SessionManager.Put(r.Context(), "userID", user.ID)

	// redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
