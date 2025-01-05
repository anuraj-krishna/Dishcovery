package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	ID            uint         `gorm:"primarykey" json:"id"`
	Name          string       `json:"name"`
	Ingredients   []Ingredient `gorm:"many2many:recipe_ingredients;" json:"-"` // Many-to-many relationship
	Steps         string       `json:"steps,omitempty"`
	Photos        string       `json:"photos,omitempty"`
	YoutubeLink   string       `json:"youtube_link,omitempty"`
	Facts         string       `json:"facts,omitempty"`
	IsVeg         bool         `json:"is_veg,omitempty"`
	Rating        float32      `json:"rating,omitempty"`
	OriginCountry string       `json:"origin_country,omitempty"`
	OriginStory   string       `json:"origin_story,omitempty"`
	CreatedAt     *time.Time   `json:"created_at,omitempty"`
	UpdatedAt     *time.Time   `json:"updated_at,omitempty"`
}

func (recipe *Recipe) GetRecipe(recipeID string) (*Recipe, error) {
	if err := db.Preload("Ingredients").First(&recipe, recipeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Recipe not found")
		}
		return nil, err
	}
	return recipe, nil
}

func (recipe *Recipe) SearchRecipe(isVeg, sortBy string, offset, limit int) ([]Recipe, error) {

	var recipes []Recipe
	queryBuilder := db.Select("id, name, rating, is_veg, photos, origin_country")

	if recipe.Name != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+recipe.Name+"%")
	}

	if isVeg != "" {
		isVegBool, err := strconv.ParseBool(isVeg)
		if err == nil {
			queryBuilder = queryBuilder.Where("is_veg = ?", isVegBool)
		}
	}

	if recipe.OriginCountry != "" {
		queryBuilder = queryBuilder.Where("origin_country ILIKE ?", "%"+recipe.OriginCountry+"%")
	}

	if sortBy == "rating" {
		queryBuilder = queryBuilder.Order("rating DESC")
	}

	queryBuilder = queryBuilder.Offset(offset).Limit(limit)
	err := queryBuilder.Find(&recipes).Error

	return recipes, err
}

func (recipe *Recipe) SearchRecipesByIngredient(ingredientId string, offset, limit int) ([]Recipe, error) {
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
		Offset(offset).Limit(limit).
		Find(&recipes).Error

	return recipes, err
}
