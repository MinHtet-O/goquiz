package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	// Declare a HTTP server using the same settings as in our main() function. srv := &http.Server{
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		//ErrorLog:     log.New(app.logger, "", 0),
	}
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
