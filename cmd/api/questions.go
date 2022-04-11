package main

import (
	"errors"
	"fmt"
	"goquiz/pkg/model"
	"goquiz/pkg/validator"
	"net/http"
)

func (app *application) getQuestionsHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()
	v := validator.New()

	categId := app.readInt(qs, "category_id", 0, v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if v.ValidateCategoryId(categId); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	category, err := app.models.CategoriesModel.Get(categId)
	if err != nil {
		fmt.Println(err.Error())
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("category", "no questions found for this category")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	questions, err := app.models.QuestionsModel.GetAll(categId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"catgory_name": category.Name, "questions": questions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
