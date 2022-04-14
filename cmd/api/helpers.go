package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"goquiz/pkg/validator"
	"net/http"
	"net/url"
	"strconv"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) (string, bool) {
	s := qs.Get(key)
	if s == "" {
		return defaultValue, false
	}
	return s, true
}

func (app *application) readCategIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	fmt.Println(params)
	id, err := strconv.ParseInt(params.ByName("categ_id"), 10, 64)
	if err != nil || id < 1 {
		fmt.Println(err.Error())
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) (int, bool) { // Extract the value from the query string.
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue, false
	}
	// Try to convert the value to an int. If this fails, add an error message to the // validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue, true
	}
	// Otherwise, return the converted integer value.
	return i, true
}
