package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.getCategoriesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.createCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/questions", app.getQuestionsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/questions", app.createQuestionHandler)
	return app.recoverPanic(app.auth(app.rateLimit(router)))
}
