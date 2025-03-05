package main

import (
	"fmt"

	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/handlers"
)

func main() {
	fmt.Println("Starting the server...")

	// Connect to Database
	database.ConnectDatabase()

	// Set up API routes
	router := handlers.SetupRouter()

	// Start the server
	router.Run(":8080") // Runs on localhost:8080
}
