package main

import "net/http"

func (app *Config) SearchRecipe(w http.ResponseWriter, r *http.Request) {
	response := app.getResponse()
	queryParams := r.URL.Query()
	rslt, err := app.Models.Recipe.SearchRecipe(queryParams)
	if err != nil {
		app.ErrorLog.Error(err)
		response.JSON(w, 400, err.Error())
		return
	}
	response.JSON(w, 200, rslt)
}

func (app *Config) SearchIngredient(w http.ResponseWriter, r *http.Request) {
	response := app.getResponse()
	queryParams := r.URL.Query()

	ingredient := app.Models.Ingredient
	ingredient.Name = queryParams.Get("name")
	rslt, err := ingredient.SearchIngredient()
	if err != nil {
		app.ErrorLog.Error(err)
		response.JSON(w, 400, err.Error())
		return
	}
	response.JSON(w, 200, rslt)
}

func (app *Config) SearchRecipesByIngredients(w http.ResponseWriter, r *http.Request) {
	response := app.getResponse()
	ingredientID := r.URL.Query().Get("ingredientsID")
	if ingredientID == "" {
		http.Error(w, "Missing 'ingredientID' query parameter", http.StatusBadRequest)
		return
	}

	recipe := app.Models.Recipe
	rslt, err := recipe.SearchRecipesByIngredient(ingredientID)
	if err != nil {
		app.ErrorLog.Error(err)
		response.JSON(w, 400, err.Error())
		return
	}
	response.JSON(w, 200, rslt)
}
