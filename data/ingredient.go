package data

import (
	"fmt"
	"time"
)

type Ingredient struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Name      string     `json:"name" gorm:"unique"` // Ingredient name
	IsVeg     bool       `json:"is_veg,omitempty"`
	Photos    string     `json:"photos"`
	Facts     string     `json:"facts"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (ingredient *Ingredient) SearchIngredient(offset, limit int) (ingredients []*Ingredient, err error) {
	name := fmt.Sprintf("%%%s%%", ingredient.Name)
	err = db.Where("ingredients.name ILIKE ?", name).
		Offset(offset).Limit(limit).
		Find(&ingredients).Error
	return ingredients, err
}
