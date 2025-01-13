package main

import (
	"dishcovery/data"
	"dishcovery/handler/dbHandler"
	"encoding/csv"
	"fmt"
	"os"
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
	for i, record := range records[1:] {
		fmt.Println(i)

		// Parse the rating from string to float32
		// rating, err := strconv.ParseFloat(record[8], 32)
		// if err != nil {
		// 	fmt.Printf("Failed to parse rating for recipe %s: %v\n", record[0], err)
		// 	continue
		// }
		// Create a new recipe object

		quantity := strings.ReplaceAll(record[2], "c.", "cup")
		// quantity = strings.ReplaceAll(quantity, "[", "")
		// quantity = strings.ReplaceAll(quantity, "]", "")
		// quantity = strings.ReplaceAll(quantity, "\"", "")
		steps := record[3]
		// steps := strings.ReplaceAll(record[3], "\"", "")
		// steps = strings.ReplaceAll(steps, "[", "")
		// steps = strings.ReplaceAll(steps, "]", "")
		recipe := data.Recipe{
			Name:     strings.Title(record[1]),
			Quantity: quantity,
			Steps:    steps,
			// WebLink:  record[4],
			IsVeg: !checkNonVeg(record[6]),
			// Name:          record[0],
			// Steps:         record[1],
			// Photos:        record[2],
			// YoutubeLink:   record[3],
			// Facts:         record[4],
			// OriginCountry: record[5],
			// OriginStory:   record[6],
			// IsVeg:         strings.ToLower(record[7]) == "true",
			// Rating:        float32(rating),
		}

		// Parse ingredients and add to the recipe3
		ingredients := strings.ReplaceAll(record[6], "[", "")
		ingredients = strings.ReplaceAll(ingredients, "]", "")

		ingredientNames := strings.Split(ingredients, ",")
		for _, ingredientName := range ingredientNames {
			ingredientName = strings.TrimSpace(ingredientName)
			ingredientName = strings.ReplaceAll(ingredientName, "\"", "")

			ingredient := data.Ingredient{Name: strings.Title(ingredientName), IsVeg: !checkNonVeg(ingredientName)}
			// caser := cases.Title(language.English)
			// caser.String()
			// Save ingredient to DB if not exists
			db.FirstOrCreate(&ingredient, data.Ingredient{Name: strings.Title(ingredientName)})
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}

		// Save recipe with its ingredients to the DB
		if err := db.Create(&recipe).Error; err != nil {
			fmt.Printf("Failed to save recipe %s: %v\n", recipe.Name, err)
		}
	}

	return nil
}

func checkNonVeg(ingredient string) bool {
	check := []string{"beef", "chicken", "bacon", "pork", "egg", "meat", "tuna",
		"hamburg", " ham ", "prawn", "mutton", "crab", "squid", "salmon"}
	for _, val := range check {
		if strings.Contains(ingredient, val) {
			return true
		} else {
			continue
		}
	}
	return false
}

func main() {
	// Initialize the database
	db = dbHandler.InitDB()

	// Load CSV data from the current directory (adjust the filename as needed)
	filename := "./../../dataset.csv" // Make sure the file is in the current directory
	err := loadCSVData(filename)
	if err != nil {
		fmt.Printf("Error loading CSV data: %v\n", err)
		return
	}

	fmt.Println("CSV data loaded and recipes saved successfully!")
}
