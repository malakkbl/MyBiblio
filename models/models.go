package models

import (
	"time"
)

// Book Model
type Book struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Title       string
	AuthorID    int
	Author      Author `gorm:"foreignKey:AuthorID"`
	Genres      string
	PublishedAt time.Time
	Price       float64
	Stock       int
}

// Author Model
type Author struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	FirstName string
	LastName  string
	Bio       string
}

// Customer Model
type Customer struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Name      string
	Email     string  `gorm:"unique"`
	Address   Address `gorm:"embedded"`
	CreatedAt time.Time
}

// Address Model
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// Order Model
type Order struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	CustomerID int
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
	TotalPrice float64
	CreatedAt  time.Time
	Status     string
}

// OrderItem Model
type OrderItem struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	OrderID  int
	BookID   int
	Book     Book `gorm:"foreignKey:BookID"`
	Quantity int
}

// SalesReport Model
type SalesReport struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	Timestamp       time.Time
	TotalRevenue    float64
	TotalOrders     int
	TopSellingBooks []BookSales `gorm:"foreignKey:ReportID"`
}

// BookSales Model
type BookSales struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	ReportID int
	BookID   int
	Book     Book `gorm:"foreignKey:BookID"`
	Quantity int
}

// SearchCriteria Model
type SearchCriteria struct {
	Titles   []string
	Authors  []string
	Genres   []string
	MinPrice float64
	MaxPrice float64
}

type User struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Role     string
}
