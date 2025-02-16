package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"um6p.ma/finalproject/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=bookstore_user password=test1234 dbname=bookstore_db port=8081 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Successfully connected to PostgreSQL!")

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

	fmt.Println("Database migrated successfully!")
}
