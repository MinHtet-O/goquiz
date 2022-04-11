package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.getCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/questions", app.getQuestionsHandler)
	return app.recoverPanic(app.rateLimit(router))
}

func (app *application) hello() {
	fmt.Println("Hello")
}
