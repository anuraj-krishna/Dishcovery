package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Config) routes() http.Handler {
	// create router
	mux := chi.NewRouter()

	// set up middleware
	mux.Use(middleware.Recoverer)

	mux.Get("/recipe", app.SearchRecipe)
	mux.Get("/recipe/{id}", app.GetRecipe)
	mux.Get("/recipe_by_ingredient", app.SearchRecipesByIngredients)
	mux.Get("/search_ingredient", app.SearchIngredient)

	return mux
}
