package main

import (
	"errors"
	"goquiz/pkg/model"
	"goquiz/pkg/validator"
	"net/http"
)

func (app *application) getQuestionsHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()
	v := validator.New()

	categId, _ := app.readInt(qs, "category_id", 0, v)
	//app.getQuestionsByCategoryId(w, r, v, categId)
	//return

	categName, _ := app.readString(qs, "category_name", "")
	//app.getQuestionsByCategoryName(w, r, v, categName)
	//return

	//
	//app.badRequestResponse(w, r, fmt.Errorf("must provide either category_id or category_name query param"))

	inputCateg := model.Category{
		ID:   categId,
		Name: categName,
	}

	questions, err := app.models.QuestionsModel.GetAll(inputCateg)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("category", "no questions found for this category name")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"questions": questions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// TODO: Postgres full text search support
// func (app *application) getQuestionsByCategoryId(w http.ResponseWriter, r *http.Request, v *validator.Validator, categId int) {
// 	if !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	if v.ValidateCategoryId(categId); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	category, err := app.models.CategoriesModel.GetByID(categId)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			v.AddError("category", "no questions found for this category id")
// 			app.failedValidationResponse(w, r, v.Errors)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	questions, err := app.models.QuestionsModel.GetAllByCategoryId(categId)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = app.writeJSON(w, http.StatusOK, envelope{"catgory_name": category.Name, "questions": questions}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
