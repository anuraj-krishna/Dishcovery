package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Recipe and Ingredient models remain the same

func initDB() {
	var err error
	dsn := "host=localhost user=admin password=localhost dbname=dishcovery_db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		os.Exit(1)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Recipe{}, &Ingredient{}); err != nil {
		fmt.Printf("Failed to migrate the database: %v\n", err)
		os.Exit(1)
	}
}

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
		// Create a new recipe object
		recipe := Recipe{
			Name:          record[0],
			Steps:         record[1],
			Photos:        record[2],
			YoutubeLink:   record[3],
			Facts:         record[4],
			OriginCountry: record[5],
			OriginStory:   record[6],
		}

		// Parse ingredients and add to the recipe
		ingredientNames := strings.Split(record[7], ",")
		for _, ingredientName := range ingredientNames {
			ingredientName = strings.TrimSpace(ingredientName)
			ingredient := Ingredient{Name: ingredientName}

			// Save ingredient to DB if not exists
			db.FirstOrCreate(&ingredient, Ingredient{Name: ingredientName})
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
	initDB() // Initialize the database

	// Load CSV data from the current directory (adjust the filename as needed)
	filename := "./recipes.csv" // Make sure the file is in the current directory
	err := loadCSVData(filename)
	if err != nil {
		fmt.Printf("Error loading CSV data: %v\n", err)
		return
	}

	fmt.Println("CSV data loaded and recipes saved successfully!")
}
