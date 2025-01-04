package main

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	Name          string       `json:"name" gorm:"index:,unique"`
	Ingredients   []Ingredient `gorm:"many2many:recipe_ingredients;" json:"ingredients"` // Many-to-many relationship
	Steps         string       `json:"steps"`
	Photos        string       `json:"photos"`
	YoutubeLink   string       `json:"youtube_link"`
	Facts         string       `json:"facts"`
	OriginCountry string       `json:"origin_country"`
	OriginStory   string       `json:"origin_story"`
}

type Ingredient struct {
	gorm.Model
	Name   string `json:"name" gorm:"index:,unique"`
	Photos string `json:"photos"`
	Facts  string `json:"facts"`
}
