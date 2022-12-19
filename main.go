package main

import (
	"os"
	"os/signal"
	"racoonapp/data"
	"racoonapp/handlers"
	"racoonapp/middleware"
	"sync"
	"syscall"

	"github.com/martijnkorbee/goracoon"
)

type application struct {
	Racoon     *goracoon.Goracoon
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
	wg         *sync.WaitGroup
}

func main() {
	app := initApplication()

	go app.listenForShutdown()

	err := app.Racoon.ListenAndServe()
	if err != nil {
		app.Racoon.ErrorLog.Panicln(err)
	}
}

func (a *application) shutdown() {
	// put any clean up tasks here

	// block untill the wait group is empty
	a.wg.Wait()

	// exit application
	os.Exit(0)
}

func (a *application) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	a.Racoon.InfoLog.Println("Received signal", s.String())

	a.shutdown()
}
