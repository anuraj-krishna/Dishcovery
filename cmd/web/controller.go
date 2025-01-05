package main

import (
	"dishcovery/data"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
)

func (app *Config) SearchRecipe(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	isVeg := queryParams.Get("is_veg")
	sortBy := queryParams.Get("sort_by")
	offset, limit := app.GetPagination(queryParams.Get("page"), queryParams.Get("limit"))

	recipe := app.Models.Recipe
	recipe.Name = queryParams.Get("name")
	recipe.OriginCountry = queryParams.Get("origin_country")

	rslt, err := recipe.SearchRecipe(isVeg, sortBy, offset, limit)
	if err != nil {
		app.ErrorLog.Error(err)
		app.FailureResponse(w, 400, "", err)
		return
	}
	app.SuccessResponse(w, 200, len(rslt), rslt)
}

func (app *Config) GetRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	rslt, err := app.Models.Recipe.GetRecipe(recipeID)
	if err != nil {
		app.ErrorLog.Error(err)
		app.FailureResponse(w, 400, "", err)
		return
	}
	app.SuccessResponse(w, 200, 0, rslt)
}

func (app *Config) SearchIngredient(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	offset, limit := app.GetPagination(queryParams.Get("page"), queryParams.Get("limit"))
	ingredient := app.Models.Ingredient
	ingredient = data.Ingredient{}
	ingredient.Name = queryParams.Get("name")
	rslt, err := ingredient.SearchIngredient(offset, limit)
	if err != nil {
		app.ErrorLog.Error(err)
		app.FailureResponse(w, 400, "", err)
		return
	}
	app.SuccessResponse(w, 200, len(rslt), rslt)
}

func (app *Config) SearchRecipesByIngredients(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ingredientID := queryParams.Get("ingredientsID")
	if ingredientID == "" {
		app.FailureResponse(w, 400, "", errors.New("missing 'ingredientsID' query parameter"))
		return
	}
	offset, limit := app.GetPagination(queryParams.Get("page"), queryParams.Get("limit"))

	recipe := app.Models.Recipe
	rslt, err := recipe.SearchRecipesByIngredient(ingredientID, offset, limit)
	if err != nil {
		app.ErrorLog.Error(err)
		app.FailureResponse(w, 400, "", err)
		return
	}
	app.SuccessResponse(w, 200, len(rslt), rslt)
}
