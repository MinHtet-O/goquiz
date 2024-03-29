package rest

import (
	"errors"
	"goquiz/service"
	"goquiz/service/validator"
	"net/http"
)

func (app *Application) getQuestionsHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()
	v := validator.New()

	categId, _ := app.readInt(qs, "category_id", 0, v)
	//cmd.getQuestionsByCategoryId(w, r, v, categId)
	//return

	categName, _ := app.readString(qs, "category_name", "")
	//cmd.getQuestionsByCategoryName(w, r, v, categName)
	//return

	//
	//cmd.badRequestResponse(w, r, fmt.Errorf("must provide either category_id or category_name query param"))

	inputCateg := service.Category{
		Id:   categId,
		Name: categName,
	}

	questions, err := app.model.QuestionsModel.GetAll(inputCateg)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
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

	for i := range questions {
		questions[i].Answer.Option++
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"questions": questions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app Application) createQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CategId     int      `json:"categ_id"`
		Text        string   `json:"text"`
		Answers     []string `json:"answers"`
		CorrectAns  int      `json:"correct_answer"`
		Codeblock   string   `json:"codeblock"`
		Explanation string   `json:"explanation"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	question := &service.Question{
		Category: service.Category{
			Id: input.CategId,
		},
		Text: input.Text,
		Answer: service.Answer{
			Explanation: input.Explanation,
			// change the ans option into array index.
			//If the correct ans is the first option, it's index is 1
			Option: service.Option(input.CorrectAns - 1),
		},
		AnsOptions: input.Answers,
	}

	v := validator.New()
	if v.ValidateQuestion(question); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	categ, err := app.model.CategoriesModel.GetByID(question.Category.Id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	id, err := app.model.QuestionsModel.Insert(categ.Id, *question)
	question.Id = id

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"question": question}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// TODO: Postgres full text search support
// func (cmd *Application) getQuestionsByCategoryId(w http.ResponseWriter, r *http.Request, v *validator.Validator, categId int) {
// 	if !v.Valid() {
// 		cmd.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	if v.ValidateCategoryId(categId); !v.Valid() {
// 		cmd.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	category, err := cmd.model.CategoriesModel.GetByID(categId)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			v.AddError("category", "no questions found for this category id")
// 			cmd.failedValidationResponse(w, r, v.Errors)
// 		default:
// 			cmd.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	questions, err := cmd.model.QuestionsModel.GetAllByCategoryId(categId)
// 	if err != nil {
// 		cmd.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = cmd.writeJSON(w, http.StatusOK, envelope{"catgory_name": category.Name, "questions": questions}, nil)
// 	if err != nil {
// 		cmd.serverErrorResponse(w, r, err)
// 	}
// }
