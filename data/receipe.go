package data

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name          string       `json:"name"`
	Ingredients   []Ingredient `gorm:"many2many:recipe_ingredients;" json:"-"` // Many-to-many relationship
	Steps         string       `json:"steps"`
	Photos        string       `json:"photos"`
	YoutubeLink   string       `json:"youtube_link"`
	Facts         string       `json:"facts"`
	IsVeg         bool         `json:"is_veg"`
	Rating        int8         `json:"rating"`
	OriginCountry string       `json:"origin_country"`
	OriginStory   string       `json:"origin_story"`
}

func (recipe *Recipe) SearchRecipe(queryParams url.Values) ([]Recipe, error) {
	query := queryParams.Get("name")
	isVeg := queryParams.Get("is_veg")
	originCountry := queryParams.Get("origin_country")
	sortBy := queryParams.Get("sort_by")
	pageStr := queryParams.Get("page")
	limitStr := queryParams.Get("limit")

	var recipes []Recipe
	queryBuilder := db.Preload("Ingredients")

	if query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+query+"%")
	}

	if isVeg != "" {
		isVegBool, err := strconv.ParseBool(isVeg)
		if err == nil {
			queryBuilder = queryBuilder.Where("is_veg = ?", isVegBool)
		}
	}

	if originCountry != "" {
		queryBuilder = queryBuilder.Where("origin_country ILIKE ?", "%"+originCountry+"%")
	}

	if sortBy == "rating" {
		queryBuilder = queryBuilder.Order("rating DESC")
	}

	page := 1
	limit := 10
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit
	queryBuilder = queryBuilder.Offset(offset).Limit(limit)

	err := queryBuilder.Find(&recipes).Error

	return recipes, err
}

func (recipe *Recipe) SearchRecipesByIngredient(ingredientId string) ([]Recipe, error) {
	var recipes []Recipe
	q := fmt.Sprintf("ingredients.id in (%s)", ingredientId)
	ingredientIds := strings.Split(ingredientId, ",")
	ingredientCount := len(ingredientIds)
	err := db.Preload("Ingredients").
		Joins("INNER JOIN recipe_ingredients ON recipes.id = recipe_ingredients.recipe_id").
		Joins("INNER JOIN ingredients ON ingredients.id = recipe_ingredients.ingredient_id").
		Where(q).
		Group("recipes.id").
		Having("COUNT(DISTINCT ingredients.id) = ?", ingredientCount).
		Find(&recipes).Error
	return recipes, err
}
