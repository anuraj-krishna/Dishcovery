package main

import (
	"net/http"
	"strconv"

	"github.com/unrolled/render"
)

func (app *Config) GetPagination(pageStr, limitStr string) (int, int) {
	page := 1
	limit := 10
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			app.ErrorLog.Error(err)
		} else if parsedPage > 0 {
			page = parsedPage
		}
	}

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			app.ErrorLog.Error(err)
		} else if parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	offset := (page - 1) * limit
	return offset, limit
}

func (app *Config) SuccessResponse(w http.ResponseWriter, status_code, count int, data any) {
	response := render.New(render.Options{StreamingJSON: true})
	r := app.Models.SuccessResponse
	r.Status = "Success"
	r.Count = count
	if status_code == 0 {
		r.StatusCode = 200
	} else {
		r.StatusCode = status_code
	}
	r.Data = data
	response.JSON(w, r.StatusCode, r)
}

func (app *Config) FailureResponse(w http.ResponseWriter, status_code int, err string, details error) {
	response := render.New(render.Options{StreamingJSON: true})
	r := app.Models.FailureResponse
	r.Status = "Failure"

	if status_code == 0 {
		r.StatusCode = 400
	} else {
		r.StatusCode = status_code
	}
	if err == "" {
		r.Error = "Somethig went wrong"
	} else {
		r.Error = err
	}
	r.Error = err
	r.ErrorDetails = details.Error()
	response.JSON(w, r.StatusCode, r)
}
