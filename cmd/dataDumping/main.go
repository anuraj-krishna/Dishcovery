package main

import (
	"dishcovery/data"
	"dishcovery/handler/dbHandler"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func loadCSVData(filename string) error {
	// Open the CSV file in the current directory
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Parse the CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	// Skip header row and process each record
	for _, record := range records[1:] {
		// Parse the rating from string to float32
		rating, err := strconv.ParseFloat(record[8], 32)
		if err != nil {
			fmt.Printf("Failed to parse rating for recipe %s: %v\n", record[0], err)
			continue
		}
		// Create a new recipe object
		recipe := data.Recipe{
			Name:          record[0],
			Steps:         record[1],
			Photos:        record[2],
			YoutubeLink:   record[3],
			Facts:         record[4],
			OriginCountry: record[5],
			OriginStory:   record[6],
			IsVeg:         strings.ToLower(record[7]) == "true",
			Rating:        float32(rating),
		}

		// Parse ingredients and add to the recipe
		ingredientNames := strings.Split(record[9], ",")
		for _, ingredientName := range ingredientNames {
			ingredientName = strings.TrimSpace(ingredientName)

			ingredient := data.Ingredient{Name: ingredientName}

			// Save ingredient to DB if not exists
			db.FirstOrCreate(&ingredient, data.Ingredient{Name: ingredientName})
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}

		// Save recipe with its ingredients to the DB
		if err := db.Create(&recipe).Error; err != nil {
			fmt.Printf("Failed to save recipe %s: %v\n", recipe.Name, err)
		}
	}

	return nil
}

func main() {
	// Initialize the database
	db = dbHandler.InitDB()

	// Load CSV data from the current directory (adjust the filename as needed)
	filename := "./recipes.csv" // Make sure the file is in the current directory
	err := loadCSVData(filename)
	if err != nil {
		fmt.Printf("Error loading CSV data: %v\n", err)
		return
	}

	fmt.Println("CSV data loaded and recipes saved successfully!")
}
