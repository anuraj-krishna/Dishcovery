package data

import (
	"fmt"

	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name   string `json:"name"` // Ingredient name
	Photos string `json:"photos"`
	Facts  string `json:"facts"`
}

func (ingredient *Ingredient) SearchIngredient() (ingredients []*Ingredient, err error) {
	name := fmt.Sprintf("%%%s%%", ingredient.Name)
	err = db.Where("ingredients.name ILIKE ?", name).
		Find(&ingredients).Error
	return ingredients, err
}
