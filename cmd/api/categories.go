package main

import (
	"net/http"
)

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: refactor as category model
	categories, err := app.models.CategoriesModel.GetCategories()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
