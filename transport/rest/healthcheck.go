package rest

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available", "system_info": map[string]string{
			"environment": app.config.Env,
			"version":     version},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
