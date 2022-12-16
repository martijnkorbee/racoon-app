package main

import (
	"fmt"
	"net/http"
	"racoonapp/data"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/martijnkorbee/goracoon/mailer"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)
	a.get("/sessions", a.Handlers.SessionsTest)
	a.get("/test-crypto", a.Handlers.TestCrypto)

	// user routes
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/login", a.Handlers.UserLogin)
	a.get("/users/logout", a.Handlers.UserLogout)

	// cache tests
	a.get("/api/cache-test", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	// mailer tests
	a.get("/mail/test", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			To:          []string{"m.korbee@numatic.nl"},
			Subject:     "test subject",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		a.Racoon.Mail.Jobs <- msg
		res := <-a.Racoon.Mail.Results
		if res.Error != nil {
			a.Racoon.ErrorLog.Println(res.Error)
		}
	})

	// db test routes
	a.get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Martijn",
			LastName:  "Korbee",
			Email:     fmt.Sprintf("%s@test.nl", a.Racoon.RandomStringGenerator(8)),
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.AddUser(u)
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "User created with ID: %d", id)
	})

	a.get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		for _, x := range users {
			fmt.Fprintf(w, "ID: %d\tFirstname: %s\tEmail: %s\n", x.ID, x.FirstName, x.Email)
		}
	})

	a.get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		user, err := a.Models.Users.GetUserByID(id)
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "ID: %d\tFirstname: %s", user.ID, user.FirstName)
	})

	a.get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		user, err := a.Models.Users.GetUserByID(id)
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		user.FirstName = a.Racoon.RandomStringGenerator(6)
		user.LastName = "changed lastname"

		validator := a.Racoon.Validator(nil)

		validator.Check(user.LastName != "", "last_name", "last name cannot be an empty string")

		if !validator.Valid() {
			fmt.Fprint(w, validator.Errors["last_name"])
			return
		}

		err = a.Models.Users.UpdateUser(*user)
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "User update: OK\n\nID: %d\tFirstname: %s", user.ID, user.FirstName)
	})

	a.get("/delete-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		err := a.Models.Users.DeleteUser(id)
		if err != nil {
			a.Racoon.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "user deleted with id: %d", id)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.Racoon.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.Racoon.Routes
}
