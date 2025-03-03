package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"um6p.ma/finalproject/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("✅ Successfully connected to PostgreSQL!")

	// AutoMigrate (creates tables if they don’t exist)
	err = DB.AutoMigrate(
		&models.Author{},
		&models.Book{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.SalesReport{},
		&models.BookSales{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("✅ Database migrated successfully!")
}
