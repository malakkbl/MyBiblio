package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"um6p.ma/finalproject/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Log the DSN (without password)
	logDSN := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)
	log.Printf("Attempting to connect to database: %s", logDSN)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)

	// Connect to PostgreSQL with detailed logging
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("✅ Successfully connected to PostgreSQL!")

	// AutoMigrate with error checking for each model
	models := []interface{}{
		&models.Author{},
		&models.Book{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.SalesReport{},
		&models.BookSales{},
		&models.User{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("❌ Failed to migrate %T: %v", model, err)
			continue
		}
		log.Printf("✅ Successfully migrated %T", model)
	}

	log.Println("✅ Database migration completed!")
}
