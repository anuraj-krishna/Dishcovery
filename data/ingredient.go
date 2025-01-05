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

func (ingredient *Ingredient) SearchIngredient(offset, limit int) (ingredients []*Ingredient, err error) {
	name := fmt.Sprintf("%%%s%%", ingredient.Name)
	err = db.Where("ingredients.name ILIKE ?", name).
		Offset(offset).Limit(limit).
		Find(&ingredients).Error
	return ingredients, err
}
