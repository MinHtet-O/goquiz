package main

import (
	"goquiz/service"
	"goquiz/service/validator"
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

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	category := &service.Category{Name: input.Name}
	v := validator.New()

	if v.ValidateCategory(category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	id, err := app.models.CategoriesModel.Insert(*category)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	category.Id = id
	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
