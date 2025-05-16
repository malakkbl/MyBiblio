package models

import (
	"time"
)

// Book Model
type Book struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"not null" validate:"required,min=1,max=200"`
	AuthorID    int       `gorm:"not null" validate:"required"`
	Author      Author    `gorm:"foreignKey:AuthorID"`
	Genres      string    `validate:"required"`
	PublishedAt time.Time `validate:"required,ltefield=now"`
	Price       float64   `validate:"required,gt=0"`
	Stock       int       `validate:"required,gte=0"`
}

// Author Model
type Author struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	FirstName string `validate:"required,min=2,max=50"`
	LastName  string `validate:"required,min=2,max=50"`
	Bio       string `validate:"max=1000"`
}

// Customer Model
type Customer struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`
	Name      string  `validate:"required,min=2,max=100"`
	Email     string  `gorm:"unique" validate:"required,email"`
	Address   Address `gorm:"embedded" validate:"required"`
	CreatedAt time.Time
}

// Address Model
type Address struct {
	Street     string `validate:"required,min=5,max=100"`
	City       string `validate:"required,min=2,max=50"`
	State      string `validate:"required,min=2,max=50"`
	PostalCode string `validate:"required,min=4,max=10"`
	Country    string `validate:"required,min=2,max=50"`
}

// Order Model
type Order struct {
	ID         int         `gorm:"primaryKey;autoIncrement"`
	CustomerID int         `validate:"required"`
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" validate:"required,min=1,dive"`
	TotalPrice float64     `validate:"required,gte=0"`
	CreatedAt  time.Time
	Status     string `validate:"required,oneof=pending processing shipped delivered cancelled"`
}

// OrderItem Model
type OrderItem struct {
	ID       int  `gorm:"primaryKey;autoIncrement"`
	OrderID  int  `validate:"required"`
	BookID   int  `validate:"required"`
	Book     Book `gorm:"foreignKey:BookID"`
	Quantity int  `validate:"required,gt=0"`
}

// SalesReport Model
type SalesReport struct {
	ID              int         `gorm:"primaryKey;autoIncrement"`
	Timestamp       time.Time   `validate:"required"`
	TotalRevenue    float64     `validate:"gte=0"`
	TotalOrders     int         `validate:"gte=0"`
	TopSellingBooks []BookSales `gorm:"foreignKey:ReportID" validate:"dive"`
}

// BookSales Model
type BookSales struct {
	ID       int  `gorm:"primaryKey;autoIncrement"`
	ReportID int  `validate:"required"`
	BookID   int  `validate:"required"`
	Book     Book `gorm:"foreignKey:BookID"`
	Quantity int  `validate:"required,gt=0"`
}

// SearchCriteria Model
type SearchCriteria struct {
	Titles   []string `validate:"dive,min=1"`
	Authors  []string `validate:"dive,min=1"`
	Genres   []string `validate:"dive,min=1"`
	MinPrice float64  `validate:"gte=0"`
	MaxPrice float64  `validate:"gtefield=MinPrice"`
}

// User Model
type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `validate:"required,min=2,max=100"`
	Email    string `gorm:"unique" validate:"required,email"`
	Password string `validate:"required,min=8,max=100"`
	Role     string `validate:"required,oneof=admin user"`
}
