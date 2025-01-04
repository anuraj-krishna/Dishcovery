package data

import "gorm.io/gorm"

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		Recipe:     Recipe{},
		Ingredient: Ingredient{},
	}
}

type Models struct {
	Recipe     Recipe
	Ingredient Ingredient
}
