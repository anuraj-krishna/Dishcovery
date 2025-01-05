package dbHandler

import (
	"dishcovery/data"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// initDB connects to Postgres and returns a pool of connections
func InitDB() *gorm.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}
	return conn
}

// connectToDB tries to connect to postgres, and backs off until a connection
// is made, or we have not connected after 10 tries
func connectToDB() *gorm.DB {
	counts := 0
	// Set up the logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                   // Enable color
		},
	)

	// dsn := os.Getenv("DATABASE_URL")
	// dsn := "host=localhost user=admin password=localhost dbname=dishcovery_db port=5432 sslmode=disable"
	dsn := "host=dpg-cttamo1opnds73cb4ht0-a.singapore-postgres.render.com user=admin password=8fE8VAP0wcipOjLwVXPXEB9jLcubjTVi dbname=dishcovery_db port=5432"

	for {
		connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Print("connected to database!")
			// Migrate the schema
			if err := connection.AutoMigrate(&data.Recipe{}, &data.Ingredient{}); err != nil {
				log.Fatalf("Failed to migrate the database: %v", err)
			}
			return connection
		}

		if counts > 3 {
			return nil
		}

		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++

		continue
	}
}
