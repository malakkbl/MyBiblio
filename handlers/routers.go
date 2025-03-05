package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns the router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Authentication Routes
	router.POST("/register", RegisterUser) // User registration
	router.POST("/login", LoginUser)       // User login

	// Books
	router.GET("/books", GetBooks)
	router.POST("/books", CreateBook)

	// Authors
	router.GET("/authors", GetAuthors)
	router.POST("/authors", CreateAuthor)

	// Customers
	router.GET("/customers", GetCustomers)
	router.POST("/customers", CreateCustomer)

	// Orders
	router.GET("/orders", GetOrders)
	router.POST("/orders", CreateOrder)

	return router
}
