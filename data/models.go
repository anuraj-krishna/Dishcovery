package data

import "gorm.io/gorm"

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		Recipe:     Recipe{},
		Ingredient: Ingredient{},
		SuccessResponse: SuccessResponse{
			Status:     "Success",
			StatusCode: 200,
		},
		FailureResponse: FailureResponse{
			Status:     "Failure",
			StatusCode: 400,
		},
	}
}

type Models struct {
	Recipe          Recipe
	Ingredient      Ingredient
	SuccessResponse SuccessResponse
	FailureResponse FailureResponse
}
