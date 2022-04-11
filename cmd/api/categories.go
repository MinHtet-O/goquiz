package main

import (
	"net/http"
)

func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	categories, err := app.models.CategoriesModel.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
